// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-30, by liasica

package pubspecmgr

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

type SegmentType int

const (
	SegmentField    SegmentType = iota // .foo
	SegmentIndex                       // [0]
	SegmentWildcard                    // [*] 或 .*
	SegmentFilter                      // ?(...)
	SegmentRoot                        // $
)

var (
	reField    = regexp.MustCompile(`^\.[a-zA-Z0-9_]+`)
	reIndex    = regexp.MustCompile(`^\[\d+]`)
	reWildcard = regexp.MustCompile(`^\[\*]|^\.\*`)
	reFilter   = regexp.MustCompile(`^\[\?\(.*?\)]`)
)

type Segment struct {
	Type  SegmentType
	Value string // 字段名 / 索引号 / 过滤器表达式
}

type PathSegments []*Segment

// ParsePathString 把路径字符串解析成 Segment 切片
func ParsePathString(s string) (PathSegments, error) {
	if s == "$" {
		return PathSegments{{Type: SegmentRoot, Value: "$"}}, nil
	}

	segments := PathSegments{{Type: SegmentRoot, Value: "$"}}
	rest := s[1:] // 去掉开头的 $

	for len(rest) > 0 {
		switch {
		case strings.HasPrefix(rest, "."):
			if m := reField.FindString(rest); m != "" {
				segments = append(segments, &Segment{Type: SegmentField, Value: m[1:]})
				rest = rest[len(m):]
				continue
			}
			if m := reWildcard.FindString(rest); m != "" {
				segments = append(segments, &Segment{Type: SegmentWildcard, Value: m})
				rest = rest[len(m):]
				continue
			}
			return nil, fmt.Errorf("%w: invalid dot path near %q", ErrInvalidPath, rest)

		case strings.HasPrefix(rest, "["):
			if m := reIndex.FindString(rest); m != "" {
				segments = append(segments, &Segment{Type: SegmentIndex, Value: m[1 : len(m)-1]})
				rest = rest[len(m):]
				continue
			}
			if m := reWildcard.FindString(rest); m != "" {
				segments = append(segments, &Segment{Type: SegmentWildcard, Value: m})
				rest = rest[len(m):]
				continue
			}
			if m := reFilter.FindString(rest); m != "" {
				segments = append(segments, &Segment{Type: SegmentFilter, Value: m})
				rest = rest[len(m):]
				continue
			}
			return nil, fmt.Errorf("%w: invalid regexp %q", ErrInvalidPath, rest)

		default:
			return nil, fmt.Errorf("%w: invalid path segment %q", ErrInvalidPath, rest)
		}
	}
	return segments, nil
}

// ParsePath 把 Path 转成 Segment 切片
func ParsePath(p *yaml.Path) (PathSegments, error) {
	return ParsePathString(p.String())
}

// Contains 判断当前路径是否包含另一路径（即另一路径是当前路径的子路径）
// 支持通配符匹配
func (sa PathSegments) Contains(sb PathSegments) bool {
	if len(sa) > len(sb) {
		return false
	}

	for i := range sa {
		if sa[i].Type == SegmentWildcard {
			// 通配符可以匹配任何
			continue
		}
		if sa[i].Value != sb[i].Value {
			return false
		}
	}
	return true
}

// PathContains 判断 a 是否包含 b（即 b 是 a 的子路径）
func PathContains(a, b *yaml.Path) bool {
	sa, err := ParsePath(a)
	if err != nil {
		return false
	}

	var sb PathSegments
	sb, err = ParsePath(b)
	if err != nil {
		return false
	}

	return sa.Contains(sb)
}
