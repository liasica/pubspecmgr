// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-31, by liasica

package subcommand

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/liasica/pubspecmgr"
)

const (
	UpgradeGroupId    = "upgrade"
	UpgradeGroupTitle = "Upgrade Commands"
)

func UpgradeCommand() *cobra.Command {
	var (
		pubspecFile string
		constraint  bool
		workers     int
	)

	cmd := &cobra.Command{
		Use:               "upgrade",
		Short:             "Upgrade pubspec dependencies",
		GroupID:           UpgradeGroupId,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(cmd *cobra.Command, args []string) {
			// check if pubspec file exists
			if _, err := os.Stat(pubspecFile); os.IsNotExist(err) {
				log.Fatalf("pubspec file %s does not exist", pubspecFile)
			}

			log.Infof("running pubspecmgr, [%d] workers on %s...", workers, pubspecFile)

			start := time.Now()
			defer func() {
				log.Infof("pubspecmgr completed in %s", time.Since(start).String())
			}()

			parsed, err := pubspecmgr.NewParser(pubspecFile).Parse()
			if err != nil {
				log.Fatalf("failed to parse pubspec file: %v", err)
			}
			upgrader := pubspecmgr.NewUpgrader(parsed, pubspecmgr.WithConstraint(constraint), pubspecmgr.WithWorkers(workers))
			upgrader.Run()
			err = parsed.Save()
			if err != nil {
				log.Fatalf("failed to save pubspec file: %v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&pubspecFile, "pubspec", "p", "./pubspec.yaml", "pubspec.yaml file path")
	cmd.Flags().BoolVarP(&constraint, "constraint", "k", false, "constraint mode, upgrade versions within the specified constraints")
	cmd.Flags().IntVarP(&workers, "workers", "w", 20, "maximum number of concurrent processes for upgrading dependencies")

	return cmd
}
