//ff:func feature=chain type=helper control=sequence
//ff:what Extracts the function body text from TS/JS lines between startLine and endLine
package chain

import "strings"

// extractTSBody extracts the function body text from lines between startLine and endLine.
func extractTSBody(lines []string, startLine, endLine int) string {
	bodyStart := startLine - 1
	bodyEnd := endLine
	if bodyStart < 0 {
		bodyStart = 0
	}
	if bodyEnd > len(lines) {
		bodyEnd = len(lines)
	}
	if bodyStart >= bodyEnd {
		return ""
	}
	return strings.Join(lines[bodyStart:bodyEnd], "\n")
}
