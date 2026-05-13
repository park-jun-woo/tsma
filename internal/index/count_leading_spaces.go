//ff:func feature=index type=helper control=iteration dimension=1
//ff:what Returns the number of leading space characters treating tabs as 4 spaces
package index

// countLeadingSpaces returns the number of leading space characters.
func countLeadingSpaces(line string) int {
	count := 0
	for _, ch := range line {
		switch ch {
		case ' ':
			count++
		case '\t':
			count += 4
		default:
			return count
		}
	}
	return count
}
