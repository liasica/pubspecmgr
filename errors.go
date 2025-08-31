// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-31, by liasica

package pubspecmgr

import "errors"

var (
	ErrPackageLine = errors.New("invalid pubspec package version line")
	ErrInvalidPath = errors.New("invalid path")
)
