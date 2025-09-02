// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-04-14, by liasica

package pubspecmgr

import (
	"os"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

// Parsed The result of parsing pubspec.yaml
type Parsed struct {
	Path     string
	File     *ast.File
	Packages []*Package
}

// Parser Parses pubspec.yaml file
type Parser struct {
	path string
}

// NewParser Create a new Parser
// filePath pubspec.yaml file path
func NewParser(filePath string) *Parser {
	return &Parser{
		path: filePath,
	}
}

// Visitor AST traversal Visitor
type Visitor struct {
	paths  []PathSegments
	marked []PathSegments

	parsed *Parsed
}

// NewVisitor Create a new Visitor
func NewVisitor(parsed *Parsed) *Visitor {
	var (
		marked []PathSegments
		paths  []PathSegments
	)

	for _, p := range GetConfig().GetMarked() {
		ps, _ := ParsePath(p)
		marked = append(marked, ps)
		top := ps[:2]

		var include bool

		for _, path := range paths {
			include = path.Contains(top)
		}

		if !include {
			paths = append(paths, top)
		}
	}

	v := &Visitor{
		paths:  paths,
		marked: marked,

		parsed: parsed,
	}

	return v
}

// Visit Implements the ast.Visitor interface
func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if !v.include(node) {
		return v
	}

	tk := node.GetToken()

	// Directly check if the node is a MappingValueNode and the token type is MappingValueType
	// This captures all mapping values, but we will filter them later based on their context
	//  and ensure they are direct string mappings
	// Example match: dependencies: pkgName: "1.0.0"
	// Non-matches: dependencies: pkgName: { git: ... } or dependencies: pkgName:
	// This is done by checking that the current token is of type MappingValueType
	// and the node is of type *ast.MappingValueNode
	//  This ensures that we only capture direct string mappings which represent package dependencies
	// and avoids capturing nested structures or non-string values.
	if tk.Type == token.MappingValueType {
		switch mv := node.(type) {
		case *ast.MappingValueNode:
			parsed, err := NewPackage(mv.Key.String(), mv.Value.String(), node)
			if err != nil {
				break
			}
			v.parsed.Packages = append(v.parsed.Packages, parsed)
		}
	}

	return v
}

// include Checks if the node's path is in the marked paths and has exactly 3 segments
func (v *Visitor) include(node ast.Node) bool {
	for _, path := range v.paths {
		ps, _ := ParsePathString(node.GetPath())

		// Only include paths that are in the marked paths and have exactly 3 segments (e.g., $.dependencies.pkgName = version)
		if path.Contains(ps) && len(ps) == 3 && !v.isMarked(ps) {
			return true
		}
	}
	return false
}

// isMarked Checks if the path segments are already marked
func (v *Visitor) isMarked(ps PathSegments) bool {
	for _, path := range v.marked {
		if path.Contains(ps) {
			return true
		}
	}
	return false
}

// Parse parses pubspec.yaml file
func (p *Parser) Parse() (parsed *Parsed, err error) {
	var b []byte
	b, err = os.ReadFile(p.path)
	if err != nil {
		return
	}

	var f *ast.File
	f, err = parser.ParseBytes(b, parser.ParseComments)
	if err != nil {
		return
	}

	parsed = &Parsed{
		File: f,
		Path: p.path,
	}

	visitor := NewVisitor(parsed)
	for _, doc := range f.Docs {
		ast.Walk(visitor, doc.Body)
	}

	return
}

// Result Returns the string representation of the parsed AST
func (p *Parsed) Result() string {
	return p.File.String()
}

// Save Saves the parsed AST back to the pubspec.yaml file
func (p *Parsed) Save() (err error) {
	str := p.Result()
	return os.WriteFile(p.Path, []byte(str), 0644)
}
