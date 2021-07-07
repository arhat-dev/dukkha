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

func (t *BaseTool) doRunTask(
	mCtx dukkha.TaskExecContext,
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
		ctx := mCtx.DeriveNew()

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

		if t.stdoutIsTty {
			stderr = utils.PrefixWriter(
				mCtx.OutputPrefix(),
				mCtx.PrefixColor(),
				mCtx.OutputColor(),
				os.Stderr,
			)
			stdout = utils.PrefixWriter(
				mCtx.OutputPrefix(),
				mCtx.PrefixColor(),
				mCtx.OutputColor(),
				os.Stdout,
			)
		} else {
			stderr = utils.PrefixWriter(
				mCtx.OutputPrefix(), nil, nil, os.Stderr,
			)
			stdout = utils.PrefixWriter(
				mCtx.OutputPrefix(), nil, nil, os.Stdout,
			)
		}

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

			if len(subSpecs) != 0 {
				err = t.doRunTask(mCtx, subSpecs, &replace)
				if err != nil {
					return fmt.Errorf("failed to run sub tasks: %w", err)
				}
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

		_, runScriptCmd, err := ctx.GetBootstrapExecSpec(cmd, false)
		if err != nil {
			return fmt.Errorf("failed to get exec spec from bootstrap config: %w", err)
		}

		output.WriteExecStart(
			t.ToolName,
			cmd,
			filepath.Base(runScriptCmd[len(runScriptCmd)-1]),
		)

		p, err := exechelper.Do(exechelper.Spec{
			Context: ctx,
			Command: runScriptCmd,
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
