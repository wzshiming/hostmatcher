package hostmatcher

import (
	"strings"
)

func newPathMatcher(path string) pathMatcher {
	if strings.Contains(path, "*") {
		return pathPatternMatcher{pattern: path}
	}
	return pathExactMatch{exact: path}
}

type pathMatcher interface {
	match(path string) bool
}

// match reports whether name matches the shell pattern.
// The pattern syntax is:
//
//	'*'       matches any sequence of non-path-separators
//	'**'      matches zero or more directories
type pathPatternMatcher struct {
	pattern string
}

func (p pathPatternMatcher) match(path string) bool {
	patternLen := len(p.pattern)
	pathLen := len(path)
	patternIndex := 0
	pathIndex := 0

	for pathIndex < pathLen && patternIndex < patternLen {
		switch p.pattern[patternIndex] {
		case '*':
			if patternIndex+1 != patternLen && p.pattern[patternIndex+1] == '*' {
				if patternIndex+2 != patternLen {
					return strings.HasSuffix(path, p.pattern[patternIndex+2:])
				}
				return true
			}

			index := strings.Index(path[pathIndex:], "/")
			if index == -1 {
				return patternLen-1 == patternIndex
			}

			pathIndex += index
			patternIndex++
		case path[pathIndex]:
			pathIndex++
			patternIndex++
			continue
		}
	}

	return pathIndex == pathLen && patternIndex == patternLen
}

type pathExactMatch struct {
	exact string
}

func (e pathExactMatch) match(path string) bool {
	return e.exact == path
}
