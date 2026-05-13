//ff:func feature=index type=helper control=sequence
//ff:what Returns true when the current line exits a TS class scope
package index

// resetTSClassContext returns true if the line should reset the class context.
func resetTSClassContext(trimmed, line string, classIndent int) bool {
	if len(trimmed) == 0 {
		return false
	}
	if countLeadingSpaces(line) > classIndent {
		return false
	}
	if tsClassPattern.MatchString(trimmed) {
		return false
	}
	if trimmed == "}" {
		return true
	}
	return false
}
