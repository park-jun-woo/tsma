//ff:func feature=index type=helper control=sequence
//ff:what Updates lastNonEmptyLine pointer when the current line is non-empty
package index

// updateLastNonEmpty sets *last to lineNum if isNonEmpty is true.
func updateLastNonEmpty(isNonEmpty bool, lineNum int, last *int) {
	if !isNonEmpty {
		return
	}
	*last = lineNum
}
