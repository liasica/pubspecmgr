// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-04-14, by liasica

package pubspecmgr

import "resty.dev/v3"

// PackageInfo represents the complete package information from pub.dev API
type PackageInfo struct {
	Name     string           `json:"name,omitempty"`
	Latest   *PackageVersion  `json:"latest,omitempty"`
	Versions []PackageVersion `json:"versions,omitempty"`
}

// PackageVersion represents a single version of a package
type PackageVersion struct {
	Version       string         `json:"version,omitempty"`
	Pubspec       PackagePubspec `json:"pubspec,omitempty"`
	ArchiveURL    string         `json:"archive_url,omitempty"`
	ArchiveSHA256 string         `json:"archive_sha256,omitempty"`
	Published     string         `json:"published,omitempty"`
}

// PackagePubspec represents the pubspec.yaml content
type PackagePubspec struct {
	Name            string         `json:"name,omitempty"`
	Version         string         `json:"version,omitempty"`
	Description     string         `json:"description,omitempty"`
	Homepage        string         `json:"homepage,omitempty"`
	Repository      string         `json:"repository,omitempty"`
	IssueTracker    string         `json:"issue_tracker,omitempty"`
	Authors         []string       `json:"authors,omitempty"`
	Author          string         `json:"author,omitempty"`
	DevDependencies map[string]any `json:"dev_dependencies,omitempty"`
	Environment     Environment    `json:"environment,omitempty"`
	Topics          []string       `json:"topics,omitempty"`
}

// Environment represents the environment constraints
type Environment struct {
	SDK *string `json:"sdk,omitempty"`
}

// GetPackageInfo fetches package information from pub.dev API
// https://pub.dev/api/packages/retry
func GetPackageInfo(pkgName string) (info *PackageInfo, err error) {
	info = new(PackageInfo)
	_, err = resty.New().R().
		SetResult(info).
		Get("https://pub.dev/api/packages/" + pkgName)
	return
}
