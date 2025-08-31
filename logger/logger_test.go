// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-31, by liasica

package logger

import (
	"testing"

	"github.com/charmbracelet/log"
)

func TestLogging(t *testing.T) {
	log.Debug("Cookie 🍪")
	log.Info("Hello World!")
	log.Error("Something went wrong")
	log.Fatal("Exiting...")
}
