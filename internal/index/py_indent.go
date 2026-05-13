//ff:func feature=index type=helper control=iteration dimension=1
//ff:what Returns the effective indentation in spaces treating tabs as 4 spaces
package index

// pyIndent returns the effective indentation in spaces (tabs = 4 spaces).
func pyIndent(s string) int {
	count := 0
	for _, ch := range s {
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
