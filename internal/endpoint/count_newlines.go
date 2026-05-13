//ff:func feature=endpoint type=helper control=sequence
//ff:what Returns the number of newline characters in a string
package endpoint

import "strings"

// countNewlines returns the number of newline characters in s.
func countNewlines(s string) int {
	return strings.Count(s, "\n")
}
