//ff:func feature=detect type=helper control=iteration dimension=1
//ff:what Finds the index of a substring within a string
package detect

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
