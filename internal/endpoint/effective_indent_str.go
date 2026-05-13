//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Computes the effective indentation width of a full line
package endpoint

// effectiveIndentStr computes the effective indentation width of a full line.
func effectiveIndentStr(line string) int {
	n := 0
	for _, ch := range line {
		if ch == ' ' {
			n++
		} else if ch == '\t' {
			n += 4
		} else {
			break
		}
	}
	return n
}
