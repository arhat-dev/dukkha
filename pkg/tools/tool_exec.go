package tools

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/output"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/pkg/exechelper"
	"arhat.dev/pkg/log"
	"github.com/fatih/color"
)

func (t *BaseTool) doRunTask(
	taskCtx *field.RenderingContext,
	outputPrefix string,
	prefixColor, outputColor *color.Color,
	execSpecs []TaskExecSpec,
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
		ctx := taskCtx.Clone()

		if es.Delay > 0 {
			_ = timer.Reset(es.Delay)

			select {
			case <-timer.C:
			case <-ctx.Context().Done():
				if !timer.Stop() {
					<-timer.C
				}

				return ctx.Context().Err()
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
			stderr = output.PrefixWriter(outputPrefix, prefixColor, outputColor, os.Stderr)
			stdout = output.PrefixWriter(outputPrefix, prefixColor, outputColor, os.Stdout)
		} else {
			stderr = output.PrefixWriter(outputPrefix, nil, nil, os.Stderr)
			stdout = output.PrefixWriter(outputPrefix, nil, nil, os.Stdout)
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
				err = t.doRunTask(taskCtx, outputPrefix, prefixColor, outputColor, subSpecs, &replace)
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
			cmd = sliceutils.NewStrings(es.Command)
		}

		_, runScriptCmd, err := t.getBootstrapExecSpec(cmd, false)
		if err != nil {
			return fmt.Errorf("failed to get exec spec from bootstrap config: %w", err)
		}

		output.WriteExecStart(
			ctx.Context(),
			t.ToolName(),
			cmd,
			filepath.Base(runScriptCmd[len(runScriptCmd)-1]),
		)

		p, err := exechelper.Do(exechelper.Spec{
			Context: ctx.Context(),
			Command: runScriptCmd,
			Env:     ctx.Values().Env,
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
