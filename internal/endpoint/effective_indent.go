//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Computes the effective indentation width from a whitespace string
package endpoint

// effectiveIndent computes the effective indentation width from a whitespace string.
func effectiveIndent(indent string) int {
	n := 0
	for _, ch := range indent {
		if ch == ' ' {
			n++
		} else if ch == '\t' {
			n += 4
		}
	}
	return n
}
