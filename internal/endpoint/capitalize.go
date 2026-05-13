//ff:func feature=endpoint type=helper control=sequence
//ff:what Returns a string with its first letter uppercased
package endpoint

import "strings"

// capitalize returns a string with its first letter uppercased.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
