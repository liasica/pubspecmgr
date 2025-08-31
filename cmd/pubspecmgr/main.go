// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-04-14, by liasica

package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/liasica/pubspecmgr"
	"github.com/liasica/pubspecmgr/cmd/pubspecmgr/subcommand"
)

func main() {
	var (
		configFile string
	)

	cmd := cobra.Command{
		Use:               "pubspecmgr",
		Short:             "pubspec version manager",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			// loading config
			pubspecmgr.LoadConfig(configFile)

			// if config file is specified, check if it exists
			// otherwise, use embedded config
			if configFile != "" {
				log.Infof("config file: %s", configFile)

				if _, err := os.Stat(configFile); os.IsNotExist(err) {
					log.Fatalf("config file %s does not exist", configFile)
				}
			}
		},
	}

	cmd.AddGroup(&cobra.Group{
		ID:    subcommand.ConfigGroupId,
		Title: subcommand.ConfigGroupTitle,
	})

	cmd.AddGroup(&cobra.Group{
		ID:    subcommand.UpgradeGroupId,
		Title: subcommand.UpgradeGroupTitle,
	})

	cmd.AddCommand(
		subcommand.Config(),
		subcommand.UpgradeCommand(),
	)

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path, default is embedded config, eg: ./pubspecmgr.yaml")

	_ = cmd.Execute()
}
