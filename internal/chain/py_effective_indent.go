//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Computes effective indentation width from a Python whitespace string
package chain

// pyEffectiveIndent computes the effective indentation width from a whitespace string.
func pyEffectiveIndent(indent string) int {
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
