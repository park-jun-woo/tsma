//ff:func feature=detect type=helper control=sequence
//ff:what Checks if a content string contains a given package string
package detect

func containsImport(content, pkg string) bool {
	return len(content) > 0 && indexOf(content, pkg) >= 0
}
