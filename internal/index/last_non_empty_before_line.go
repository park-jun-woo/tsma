//ff:func feature=index type=helper control=sequence
//ff:what Returns the last non-empty line number before the current declaration line
package index

// lastNonEmptyBeforeLine returns the last non-empty line number before the
// current declaration line. If lastNonEmpty is before the current line, use it;
// otherwise fall back to one line before current.
func lastNonEmptyBeforeLine(currentLine, lastNonEmpty int) int {
	if lastNonEmpty < currentLine {
		return lastNonEmpty
	}
	if currentLine > 1 {
		return currentLine - 1
	}
	return currentLine
}
