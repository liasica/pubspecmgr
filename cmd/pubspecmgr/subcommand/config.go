// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-31, by liasica

package subcommand

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"

	"github.com/liasica/pubspecmgr"
)

const (
	ConfigGroupId    = "config"
	ConfigGroupTitle = "Configuration Commands"
)

func Config() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "config",
		Short:             "Manage configuration",
		GroupID:           ConfigGroupId,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	cmd.AddCommand(
		configPrint(),
		configCreate(),
	)

	return cmd
}

func configPrint() *cobra.Command {
	return &cobra.Command{
		Use:               "print",
		Short:             "Print current configuration",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(cmd *cobra.Command, args []string) {
			b, err := yaml.Marshal(pubspecmgr.GetConfig())
			if err != nil {
				log.Fatalf("failed to marshal config: %v", err)
			}
			log.Infof("current configuration yaml:\n```yaml\n%s\n```", strings.Trim(string(b), "\n"))
		},
	}
}

func configCreate() *cobra.Command {
	return &cobra.Command{
		Use:               "create",
		Short:             "Create a new configuration file on current directory",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("creating new configuration file...")
			b, err := yaml.Marshal(pubspecmgr.GetConfig())
			if err != nil {
				log.Fatalf("failed to marshal config: %v", err)
			}
			err = os.WriteFile("pubspecmgr.yaml", b, 0644)
			if err != nil {
				log.Fatalf("failed to write config file: %v", err)
			}
			log.Infof("new configuration file created at ./pubspecmgr.yaml")
		},
	}
}
