//ff:func feature=index type=helper control=sequence
//ff:what Returns true if a line opens a Rust impl block body
package index

import "strings"

// isRsImplLine reports whether the trimmed line matches an impl declaration
// whose body opens on the same line (contains an opening brace).
func isRsImplLine(trimmed, line string) bool {
	return rsImplPattern.MatchString(trimmed) && strings.Contains(line, "{")
}
