package tools

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"arhat.dev/pkg/exechelper"
	"arhat.dev/pkg/log"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/output"
	"arhat.dev/dukkha/pkg/utils"
)

// nolint:gocyclo
func doRun(
	ctx dukkha.TaskExecContext,
	execSpecs []dukkha.TaskExecSpec,
	_replaceEntries *dukkha.ReplaceEntries,
) error {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}

	var replace dukkha.ReplaceEntries
	if _replaceEntries != nil {
		replace = *_replaceEntries
	} else {
		replace = make(dukkha.ReplaceEntries)
	}

	for _, es := range execSpecs {
		if es.Delay > 0 {
			_ = timer.Reset(es.Delay)

			select {
			case <-timer.C:
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}

				return ctx.Err()
			}
		}

		var (
			stdin          io.Reader
			stdout, stderr io.Writer
		)

		if es.Stdin != nil {
			stdin = es.Stdin
		} else {
			stdin = os.Stdin
		}

		stderr = utils.PrefixWriter(
			ctx.OutputPrefix(), ctx.ColorOutput(),
			ctx.PrefixColor(), ctx.OutputColor(),
			os.Stderr,
		)
		stdout = utils.PrefixWriter(
			ctx.OutputPrefix(), ctx.ColorOutput(),
			ctx.PrefixColor(), ctx.OutputColor(),
			os.Stdout,
		)

		var (
			stdoutBuf *bytes.Buffer
			stderrBuf *bytes.Buffer
		)
		if len(es.StdoutAsReplace) != 0 {
			stdoutBuf = &bytes.Buffer{}
			stdout = io.MultiWriter(stdout, stdoutBuf)
		}

		if len(es.StderrAsReplace) != 0 {
			stderrBuf = &bytes.Buffer{}
			stderr = io.MultiWriter(stderr, stderrBuf)
		}

		setReplaceEntry := func(err error) {
			if stdoutBuf != nil {
				stdoutValue := stdoutBuf.Bytes()
				if es.FixStdoutValueForReplace != nil {
					stdoutValue = es.FixStdoutValueForReplace(stdoutValue)
				}

				replace[es.StdoutAsReplace] = &dukkha.ReplaceEntry{
					Data: stdoutValue,
					Err:  err,
				}
			}

			if stderrBuf != nil {
				stderrValue := stderrBuf.Bytes()
				if es.FixStdoutValueForReplace != nil {
					stderrValue = es.FixStderrValueForReplace(stderrValue)
				}

				replace[es.StderrAsReplace] = &dukkha.ReplaceEntry{
					Data: stderrValue,
					Err:  err,
				}
			}
		}

		// alter exec func can generate sub exec specs
		if es.AlterExecFunc != nil {
			subSpecs, err := es.AlterExecFunc(replace, stdin, stdout, stderr)
			setReplaceEntry(err)
			if err != nil {
				return fmt.Errorf("failed to execute alter exec func: %w", err)
			}

			switch t := subSpecs.(type) {
			case []dukkha.TaskExecSpec:
				err = doRun(ctx, t, &replace)
			case *TaskExecRequest:
				err = RunTask(t)
			case nil:
				// nothing to do
			default:
				// TODO: log error instead of panic?
				panic(fmt.Errorf("unexpected sub specs type: %T", t))
			}

			if err != nil {
				return fmt.Errorf("failed to run sub tasks: %w", err)
			}

			continue
		}

		cmd := es.Command
		if len(replace) != 0 {
			pairs := make([]string, 2*len(replace))
			i := 0
			for toReplace, newValue := range replace {
				pairs[i], pairs[i+1] = toReplace, string(newValue.Data)
				i += 2
			}

			replacer := strings.NewReplacer(pairs...)

			// replace placeholders in cmd
			cmd = make([]string, 0, len(es.Command))
			for _, origCmdPart := range es.Command {
				cmd = append(cmd, replacer.Replace(origCmdPart))
			}

			// replace placeholders in env
			for _, origEnvPart := range es.EnvOverride {
				ctx.AddEnv(true, dukkha.EnvEntry{
					Name:  replacer.Replace(origEnvPart.Name),
					Value: replacer.Replace(origEnvPart.Value),
				})
			}

			for _, origEnvPart := range es.EnvSuggest {
				ctx.AddEnv(false, dukkha.EnvEntry{
					Name:  replacer.Replace(origEnvPart.Name),
					Value: replacer.Replace(origEnvPart.Value),
				})
			}
		} else {
			ctx.AddEnv(true, es.EnvOverride...)
			ctx.AddEnv(false, es.EnvSuggest...)
		}

		var err error
		if es.UseShell {
			var shellCmd []string
			if es.ShellName != "" {
				// using embedded shell
				sh, ok := ctx.GetShell(es.ShellName)
				if !ok {
					return fmt.Errorf("shell %q not found", es.ShellName)
				}

				_, shellCmd, err = sh.GetExecSpec(cmd, false)
				if err != nil {
					return fmt.Errorf("failed to get exec spec for shell %q: %w", es.ShellName, err)
				}
			}

			output.WriteExecStart(
				ctx.PrefixColor(),
				ctx.CurrentTool(), cmd,
				filepath.Base(shellCmd[len(shellCmd)-1])[:7],
			)

			cmd = shellCmd
		} else {
			output.WriteExecStart(ctx.PrefixColor(), ctx.CurrentTool(), cmd, "")
		}

		p, err := exechelper.Do(exechelper.Spec{
			Context: ctx,
			Command: cmd,
			Env:     ctx.Env(),
			Dir:     es.Chdir,

			Stdin: stdin,

			Stdout: stdout,
			Stderr: stderr,
		})
		if err != nil {
			setReplaceEntry(err)
			if !es.IgnoreError {
				return fmt.Errorf("failed to prepare command [ %s ]: %w", strings.Join(cmd, " "), err)
			}

			// TODO: log error in detail
			log.Log.I("error ignored", log.Error(err))

			continue
		}

		_, err = p.Wait()
		setReplaceEntry(err)

		if err != nil {
			if !es.IgnoreError {
				return fmt.Errorf("command exited with error: %w", err)
			}

			// TODO: log error in detail
			log.Log.I("error ignored", log.Error(err))
			continue
		}
	}

	return nil
}
