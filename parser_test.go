// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-30, by liasica

package pubspecmgr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	LoadConfig()

	p := NewParser("testdata/pubspec.yaml")
	parsed, err := p.Parse()
	require.NoError(t, err)
	t.Logf("Parsed %d packages", len(parsed.Packages))
}
