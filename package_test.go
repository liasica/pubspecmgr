// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-31, by liasica

package pubspecmgr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPackageFromString(t *testing.T) {
	p, err := NewPackageFromString("\n  intl: ^0.20.2", nil)
	require.NoError(t, err)
	require.Equal(t, "intl", p.Name)
	require.Equal(t, "^0.20.2", p.Version)

	p, err = NewPackageFromString("  flutter_inappwebview: 6.1.5", nil)
	require.NoError(t, err)
	require.Equal(t, "flutter_inappwebview", p.Name)
	require.Equal(t, "6.1.5", p.Version)

	p, err = NewPackageFromString("google_maps_flutter: 2.12.3", nil)
	require.NoError(t, err)
	require.Equal(t, "google_maps_flutter", p.Name)
	require.Equal(t, "2.12.3", p.Version)

	p, err = NewPackageFromString("  cupertino_icons: ^1.0.2 # some comment", nil)
	require.NoError(t, err)
	require.Equal(t, "cupertino_icons", p.Name)
	require.Equal(t, "^1.0.2", p.Version)

	p, err = NewPackageFromString("  some_package: # no version", nil)
	require.Error(t, err)
	require.ErrorAs(t, err, &ErrPackageLine)
	require.Nil(t, p)

	p, err = NewPackageFromString("invalid_line_without_colon", nil)
	require.Error(t, err)
	require.Nil(t, p)
}
