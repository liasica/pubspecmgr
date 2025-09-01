// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-09-01, by liasica

package pubspecmgr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpgraderRun(t *testing.T) {
	LoadConfig()

	parsed, err := NewParser("testdata/pubspec.yaml").Parse()
	require.NoError(t, err)
	upgrader := NewUpgrader(parsed, WithConstraint(false), WithWorkers(20))
	upgrader.Run()
	result := parsed.Result()
	fmt.Println(result)
}
