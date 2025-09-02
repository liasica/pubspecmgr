// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-30, by liasica

package pubspecmgr

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/charmbracelet/log"
	"github.com/goccy/go-yaml/ast"
)

var packageRegexp = regexp.MustCompile(`^(\S+)?:[| ]([\^|]?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?)`)

// Package represents a Dart/Flutter package with its name and version
//
// Semver: https://semver.org/lang/zh-CN/
type Package struct {
	node    ast.Node
	Name    string
	Version string

	constraints *semver.Constraints
	info        *PackageInfo
}

// NewPackage creates a new Package instance
func NewPackage(name, version string, node ast.Node) (*Package, error) {
	constraints, err := semver.NewConstraint(version)
	if err != nil {
		return nil, err
	}

	return &Package{
		Name:        name,
		Version:     version,
		node:        node,
		constraints: constraints,
	}, nil
}

// NewPackageFromString Creates a new Package instance from a string
// The string should be in the format "package: version"
// Returns nil if the string does not match the expected format
func NewPackageFromString(str string, node ast.Node) (*Package, error) {
	matches := packageRegexp.FindStringSubmatch(strings.TrimSpace(str))
	if len(matches) < 3 {
		return nil, fmt.Errorf("%w: %s", ErrPackageLine, str)
	}

	name := strings.TrimSpace(matches[1])
	version := strings.TrimSpace(matches[2])

	return NewPackage(name, version, node)
}

// GetLatest Fetches the latest version of the package from the pub.dev API
func (pkg *Package) GetLatest(constraint bool) (newVersion string) {
	newVersion = pkg.Version

	var prefix string
	if strings.HasPrefix(pkg.Version, "^") {
		prefix = "^"
	}

	var err error
	pkg.info, err = GetPackageInfo(pkg.Name)
	if err != nil {
		log.Error("failed to get package info: %v", err)
		return
	}

	// If constraint is true, find the latest version that satisfies the constraint
	if constraint {
		for i := len(pkg.info.Versions) - 1; i >= 0; i-- {
			pver := pkg.info.Versions[i]

			// Parse the version string
			var ver *semver.Version
			ver, err = semver.NewVersion(pver.Version)
			if err != nil {
				log.Error("failed to parse package version: %v", err)
				return
			}

			// Check if the version satisfies the constraint
			if pkg.constraints.Check(ver) {
				newVersion = prefix + pver.Version
				break
			}
		}
	} else {
		// If no constraint, just take the latest version
		if pkg.info.Latest != nil {
			newVersion = prefix + pkg.info.Latest.Version
		}
	}

	if newVersion != pkg.Version {
		log.Infof("latest version of package %s is %s, current: %s", pkg.Name, newVersion, pkg.Version)
	}

	return
}
