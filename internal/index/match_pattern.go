//ff:func feature=index type=util control=sequence
//ff:what 단일 .tsmignore 패턴과 경로를 매칭
package index

import (
	"path/filepath"
	"strings"
)

func matchPattern(path string, name string, isDir bool, pattern string) bool {
	if strings.HasSuffix(pattern, "/") {
		dirPattern := strings.TrimSuffix(pattern, "/")
		if !isDir {
			return false
		}
		if name == dirPattern {
			return true
		}
		return strings.HasSuffix(path, dirPattern) || strings.Contains("/"+path+"/", "/"+dirPattern+"/")
	}
	if strings.Contains(pattern, "/") {
		clean := strings.TrimPrefix(path, "./")
		matched, _ := filepath.Match(pattern, clean)
		return matched
	}
	matched, _ := filepath.Match(pattern, name)
	return matched
}
