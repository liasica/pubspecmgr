// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-31, by liasica

package g

import "fmt"

var (
	Version    = "1.0.5"
	Commit     = "d49cdcb"
	CommitDate = "2026-04-15T19:24:25+08:00"
)

func GetVersion() string {
	return fmt.Sprintf("%s-%s (%s)", Version, Commit, CommitDate)
}
