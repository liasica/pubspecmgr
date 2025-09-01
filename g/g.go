// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-31, by liasica

package g

import "fmt"

var (
	Version    = "1.0.0"
	Commit     = "none"
	CommitDate = "unknown"
)

func GetVersion() string {
	return fmt.Sprintf("%s-%s (%s)", Version, Commit, CommitDate)
}
