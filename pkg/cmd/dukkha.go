/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"arhat.dev/pkg/log"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/shell_file"
	"arhat.dev/dukkha/pkg/renderer/template"
	"arhat.dev/dukkha/pkg/renderer/template_file"

	_ "embed" // go:embed to embed default config file
)

func NewRootCmd() *cobra.Command {
	var (
		appCtx      context.Context
		configPaths []string
		config      = conf.NewConfig()
	)

	rootCmd := &cobra.Command{
		Use:           "dukkha",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "version" {
				return nil
			}

			err := log.SetDefaultLogger(log.ConfigSet{config.Log})
			if err != nil {
				return err
			}

			err = readConfig(
				configPaths,
				cmd.PersistentFlags().Changed("config"),
				&config,
			)
			if err != nil {
				return fmt.Errorf("failed to read config: %w", err)
			}

			appCtx = context.Background()

			err = config.Bootstrap.Resolve()
			if err != nil {
				return fmt.Errorf("failed to resolve bootstrap config: %w", err)
			}

			// bootstrap config was resolved when unmarshaling
			if config.Bootstrap.Shell == "" {
				return fmt.Errorf("unable to get a shell name, please set bootstrap.shell manually")
			}

			// create a renderer manager with essential renderers
			mgr := renderer.NewManager()
			err = multierr.Append(
				err,
				mgr.Add(
					&shell.Config{ExecFunc: config.Bootstrap.Exec},
					shell.DefaultName,
				),
			)
			err = multierr.Append(
				err,
				mgr.Add(
					&shell_file.Config{ExecFunc: config.Bootstrap.Exec},
					shell_file.DefaultName,
				),
			)
			err = multierr.Append(
				err,
				mgr.Add(&template.Config{}, template.DefaultName),
			)
			err = multierr.Append(
				err,
				mgr.Add(&template_file.Config{}, template_file.DefaultName),
			)
			err = multierr.Append(
				err,
				mgr.Add(&file.Config{}, file.DefaultName),
			)
			if err != nil {
				return fmt.Errorf("failed to create essential renderers: %w", err)
			}

			appCtx = renderer.WithManager(appCtx, mgr)

			// ensure all top-level config resolved
			err = config.Resolve(appCtx, mgr.Render, 1)
			if err != nil {
				return fmt.Errorf("failed to resolve config: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(appCtx, config)
		},
	}

	globalFlags := rootCmd.PersistentFlags()

	globalFlags.StringSliceVarP(&configPaths, "config", "c", []string{".dukkha", ".dukkha.yaml"}, "")

	// logging
	globalFlags.StringVarP(&config.Log.Level, "log.level", "v", "info",
		"log level, one of [verbose, debug, info, error, silent]")
	globalFlags.StringVar(&config.Log.Format, "log.format", "console",
		"log output format, one of [console, json]")
	globalFlags.StringVar(&config.Log.File, "log.file", "stderr",
		"log output to this file")

	return rootCmd
}

func run(appCtx context.Context, config *conf.Config) error {
	logger := log.Log.WithName("app")

	logger.I("application configured",
		log.Any("bootstrap", config.Bootstrap),
		log.Any("shell", config.Shell),
	)

	_ = appCtx
	return nil
}

// do not use strict unmarshal when reading config, tasks are dynamic
func readConfig(configPaths []string, failOnFileNotFoundError bool, mergedConfig **conf.Config) error {
	readAndMergeConfigFile := func(path string) error {
		configBytes, err2 := os.ReadFile(path)
		if err2 != nil {
			return fmt.Errorf("failed to read config file %q: %w", path, err2)
		}

		current := conf.NewConfig()
		err2 = yaml.Unmarshal(configBytes, &current)
		if err2 != nil {
			return fmt.Errorf("failed to unmarshal config file %q: %w", path, err2)
		}

		log.Log.V("config unmarshaled", log.String("file", path), log.Any("config", current))

		// TODO: merge into mergedConfig
		_ = mergedConfig

		return err2
	}

	for _, path := range configPaths {
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				if !failOnFileNotFoundError {
					continue
				}
			}

			return err
		}

		if !info.IsDir() {
			err = readAndMergeConfigFile(path)
			if err != nil {
				return err
			}

			continue
		}

		err = fs.WalkDir(os.DirFS(path), ".", func(pathInDir string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			switch filepath.Ext(pathInDir) {
			case ".yml", ".yaml":
			default:
				return nil
			}

			return readAndMergeConfigFile(filepath.Join(path, pathInDir))
		})

		if err != nil {
			return err
		}
	}

	return nil
}
