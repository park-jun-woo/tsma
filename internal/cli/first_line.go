//ff:func feature=cli type=util control=iteration dimension=1
//ff:what Returns the first line of a string for compact error logging
package cli

// firstLine returns the first line of s, used to keep skipped-package failure
// logs to a single line during the batch coverage scan.
func firstLine(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			return s[:i]
		}
	}
	return s
}
