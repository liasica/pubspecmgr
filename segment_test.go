// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-30, by liasica

package pubspecmgr

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
)

func TestPathContains(t *testing.T) {
	p1, err := yaml.PathString("$.a.b.c")
	require.NoError(t, err)

	var s1 PathSegments
	s1, err = ParsePath(p1)
	require.NoError(t, err)

	var p2 *yaml.Path
	p2, err = yaml.PathString("$.a.b")
	require.NoError(t, err)

	var s2 PathSegments
	s2, err = ParsePath(p2)
	require.NoError(t, err)

	require.True(t, s2.Contains(s1))
	require.True(t, PathContains(p2, p1))

	var v1 *semver.Version
	v1, err = semver.NewVersion("2.5.0")
	require.NoError(t, err)

	var c1 *semver.Constraints
	c1, err = semver.NewConstraint(">=2.3.0 <3.0.0")
	require.NoError(t, err)

	require.True(t, c1.Check(v1))

	// if !p1.Contains(p2) {
	// 	t.Fatalf("expected %q to contain %q", p1, p2)
	// }
	// if p2.Contains(p1) {
	// 	t.Fatalf("expected %q to not contain %q", p2, p1)
	// }
}
