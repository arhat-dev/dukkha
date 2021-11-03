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

	var ansiTranslationExitSig chan struct{}

	notifyLastANSITranslationExit := func() {
		if ansiTranslationExitSig != nil {
			select {
			case <-ansiTranslationExitSig:
				// signaled
			default:
				// not signaled
				close(ansiTranslationExitSig)
			}
		}
	}

	defer notifyLastANSITranslationExit()

	for _, es := range execSpecs {
		notifyLastANSITranslationExit()

		if es.Delay > 0 {
			_ = timer.Reset(es.Delay)

			select {
			case <-timer.C:
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}

				ctx.SetState(dukkha.TaskExecCanceled)
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

		if !ctx.TranslateANSIStream() {
			stdout = os.Stdout
			stderr = os.Stderr
		} else {
			stdoutW := utils.NewANSIWriter(
				os.Stdout, ctx.RetainANSIStyle(),
			)

			stdout = stdoutW
			stderr = stdoutW

			ansiTranslationExitSig = make(chan struct{})

			go func() {
				// TODO: make flush interval customizable
				ticker := time.NewTicker(2 * time.Second)

				defer func() {
					ticker.Stop()
					_, err := stdoutW.Flush()
					if err != nil {
						log.Log.I(
							"failed to flush translated plain text data to stdout when closing",
							log.Error(err),
						)
					}
				}()

				for {
					select {
					case <-ticker.C:
						_, err := stdoutW.Flush()
						if err != nil {
							log.Log.I(
								"failed to flush translated plain text data to stdout",
								log.Error(err),
							)
							return
						}
					case <-ansiTranslationExitSig:
						return
					}
				}
			}()
		}

		stdout = utils.TermWriter(
			ctx.OutputPrefix(), ctx.ColorOutput(),
			ctx.PrefixColor(), ctx.OutputColor(),
			stdout,
		)

		stderr = utils.TermWriter(
			ctx.OutputPrefix(), ctx.ColorOutput(),
			ctx.PrefixColor(), ctx.OutputColor(),
			stderr,
		)

		var (
			stdoutBuf *bytes.Buffer
			stderrBuf *bytes.Buffer
		)
		if len(es.StdoutAsReplace) != 0 {
			stdoutBuf = &bytes.Buffer{}

			if es.ShowStdout {
				stdout = io.MultiWriter(stdout, stdoutBuf)
			} else {
				stdout = stdoutBuf
			}
		}

		if len(es.StderrAsReplace) != 0 {
			stderrBuf = &bytes.Buffer{}

			if es.ShowStderr {
				stderr = io.MultiWriter(stderr, stderrBuf)
			} else {
				stderr = stderrBuf
			}
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
				if es.FixStderrValueForReplace != nil {
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
			ctx.SetState(dukkha.TaskExecWorking)

			subSpecs, err := es.AlterExecFunc(replace, stdin, stdout, stderr)
			setReplaceEntry(err)
			if err != nil {
				ctx.SetState(dukkha.TaskExecFailed)
				return err
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
				ctx.SetState(dukkha.TaskExecFailed)
				return err
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
				ctx.AddEnv(true, &dukkha.EnvEntry{
					Name:  replacer.Replace(origEnvPart.Name),
					Value: replacer.Replace(origEnvPart.Value),
				})
			}

			for _, origEnvPart := range es.EnvSuggest {
				ctx.AddEnv(false, &dukkha.EnvEntry{
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
					ctx.SetState(dukkha.TaskExecFailed)
					return fmt.Errorf("shell %q not found", es.ShellName)
				}

				_, shellCmd, err = sh.GetExecSpec(cmd, false)
				if err != nil {
					ctx.SetState(dukkha.TaskExecFailed)
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

		ctx.SetState(dukkha.TaskExecWorking)
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
			ctx.SetState(dukkha.TaskExecFailed)
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
			ctx.SetState(dukkha.TaskExecFailed)
			if !es.IgnoreError {
				return fmt.Errorf("command exited with error: %w", err)
			}

			// TODO: log error in detail
			log.Log.I("error ignored", log.Error(err))
			continue
		}

		ctx.SetState(dukkha.TaskExecSucceeded)
	}

	return nil
}
