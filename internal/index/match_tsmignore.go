//ff:func feature=index type=util control=iteration dimension=1
//ff:what 경로가 .tsmignore 패턴에 매칭되는지 판별
package index

// MatchTsmIgnore returns true if the given path matches any of the ignore patterns.
func MatchTsmIgnore(path string, name string, isDir bool, patterns []string) bool {
	for _, pattern := range patterns {
		if matchPattern(path, name, isDir, pattern) {
			return true
		}
	}
	return false
}
