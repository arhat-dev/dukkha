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
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/utils"
)

// nolint:gocyclo
func doRun(
	ctx dukkha.TaskExecContext,
	execSpecs []dukkha.TaskExecSpec,
	_replaceEntries *map[string][]byte,
) error {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}

	var replace map[string][]byte
	if _replaceEntries != nil {
		replace = *_replaceEntries
	} else {
		replace = make(map[string][]byte)
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
			stdin = io.MultiReader(es.Stdin, os.Stdin)
		} else {
			stdin = os.Stdin
		}

		stderr = utils.PrefixWriter(
			ctx.OutputPrefix(),
			ctx.ColorOutput(),
			ctx.PrefixColor(),
			ctx.OutputColor(),
			os.Stderr,
		)
		stdout = utils.PrefixWriter(
			ctx.OutputPrefix(),
			ctx.ColorOutput(),
			ctx.PrefixColor(),
			ctx.OutputColor(),
			os.Stdout,
		)

		var buf *bytes.Buffer
		if len(es.OutputAsReplace) != 0 {
			buf = &bytes.Buffer{}

			stdout = io.MultiWriter(stdout, buf)
		}

		// alter exec func can generate sub exec specs
		if es.AlterExecFunc != nil {
			subSpecs, err := es.AlterExecFunc(replace, stdin, stdout, stderr)
			if err != nil {
				return fmt.Errorf("failed to execute alter exec func: %w", err)
			}

			if buf != nil {
				newValue := buf.Bytes()
				if es.FixOutputForReplace != nil {
					newValue = es.FixOutputForReplace(newValue)
				}

				replace[es.OutputAsReplace] = newValue
			}

			switch t := subSpecs.(type) {
			case []dukkha.TaskExecSpec:
				err = doRun(ctx, t, &replace)
			case *CompleteTaskExecSpecs:
				err = RunTask(t.Context, t.Tool, t.Task)
			case nil:
			default:
				// TODO: log error instead of panic?
				panic(fmt.Errorf("unexpected sub specs type: %T", t))
			}

			if err != nil {
				return fmt.Errorf("failed to run sub tasks: %w", err)
			}

			continue
		}

		var cmd []string
		if len(replace) != 0 {
			pairs := make([]string, 2*len(replace))
			i := 0
			for toReplace, newValue := range replace {
				pairs[i], pairs[i+1] = toReplace, string(newValue)
				i += 2
			}

			replacer := strings.NewReplacer(pairs...)
			for _, rawEnvPart := range es.Env {
				ctx.AddEnv(replacer.Replace(rawEnvPart))
			}

			for _, rawCmdPart := range es.Command {
				cmd = append(cmd, replacer.Replace(rawCmdPart))
			}
		} else {
			ctx.AddEnv(es.Env...)
			cmd = sliceutils.NewStrings(es.Command)
		}

		var err error
		if es.UseShell {
			var shellCmd []string
			if es.ShellName == "bootstrap" {
				_, shellCmd, err = ctx.GetBootstrapExecSpec(cmd, false)
				if err != nil {
					return fmt.Errorf("failed to get exec spec for bootstrap shell: %w", err)
				}
			} else {
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
			if !es.IgnoreError {
				return fmt.Errorf("failed to prepare command [ %s ]: %w", strings.Join(cmd, " "), err)
			}

			// TODO: log error in detail
			log.Log.I("error ignored", log.Error(err))

			delete(replace, es.OutputAsReplace)

			continue
		}

		_, err = p.Wait()
		if err != nil {
			if !es.IgnoreError {
				return fmt.Errorf("command exited with error: %w", err)
			}

			// TODO: log error in detail
			log.Log.I("error ignored", log.Error(err))

			delete(replace, es.OutputAsReplace)

			continue
		}

		if buf != nil {
			newValue := buf.Bytes()
			if es.FixOutputForReplace != nil {
				newValue = es.FixOutputForReplace(newValue)
			}

			replace[es.OutputAsReplace] = newValue
		}
	}

	return nil
}
