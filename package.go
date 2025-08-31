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

var packageRegexp = regexp.MustCompile(`(\S+)?:[| ]([\^|]?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?)`)

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
func NewPackage(name, version string, node ast.Node) *Package {
	constraints, _ := semver.NewConstraint(version)

	return &Package{
		Name:        name,
		Version:     version,
		node:        node,
		constraints: constraints,
	}
}

// NewPackageFromString Creates a new Package instance from a string
// The string should be in the format "package: version"
// Returns nil if the string does not match the expected format
func NewPackageFromString(str string, node ast.Node) (*Package, error) {
	matches := packageRegexp.FindStringSubmatch(str)
	if len(matches) < 3 {
		return nil, fmt.Errorf("%w: %s", ErrPackageLine, str)
	}

	name := strings.TrimSpace(matches[1])
	version := strings.TrimSpace(matches[2])

	return NewPackage(name, version, node), nil
}

// GetLatest Fetches the latest version of the package from the pub.dev API
func (pkg *Package) GetLatest() (newVersion string) {
	newVersion = pkg.Version

	var err error
	pkg.info, err = GetPackageInfo(pkg.Name)
	if err != nil {
		log.Error("failed to get package info: %v", err)
		return
	}

	log.Infof("latest version of package %s is %s, current: %s", pkg.Name, pkg.info.Latest.Version, pkg.Version)

	return pkg.info.Latest.Version
}
