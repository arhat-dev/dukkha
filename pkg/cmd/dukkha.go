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
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/conf"

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
				configPaths, cmd.PersistentFlags().Changed("config"),
				config,
			)
			if err != nil {
				return fmt.Errorf("failed to read config: %w", err)
			}

			appCtx = context.Background()
			appCtx, err = config.Resolve(appCtx)
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
		// log.Any("tools", config.Tools),
	)

	_ = appCtx
	return nil
}

//go:embed default.yaml
var defaultConfigBytes []byte

// do not use strict unmarshal when reading config, tasks are dynamic
func readConfig(configPaths []string, failOnFileNotFoundError bool, mergedConfig *conf.Config) error {
	err := yaml.Unmarshal(defaultConfigBytes, mergedConfig)
	if err != nil {
		return fmt.Errorf("invalid default config: %w", err)
	}

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

		// err2 = mergo.Merge(&mergedConfig, current, mergo.WithOverride)
		// if err2 != nil {
		// 	return fmt.Errorf("failed to merge config file %q: %w", path, err2)
		// }

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
