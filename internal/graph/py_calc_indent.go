//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Returns the indentation level of a line treating tabs as 4 spaces
package graph

// pyCalcIndent returns the indentation level of a line (tabs = 4 spaces).
func pyCalcIndent(line string) int {
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
