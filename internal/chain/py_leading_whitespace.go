//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Extracts the leading whitespace prefix from a line of Python code
package chain

// pyLeadingWhitespace extracts leading whitespace from a line.
func pyLeadingWhitespace(line string) string {
	for i, ch := range line {
		if ch != ' ' && ch != '\t' {
			return line[:i]
		}
	}
	return line
}
