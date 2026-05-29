//ff:func feature=index type=helper control=iteration dimension=1
//ff:what Returns the index of the closing quote for a string or char literal, honoring escapes and Rust lifetimes
package index

// skipQuoted scans a string or char literal that opens at line[start] and
// returns the index of its closing quote. Backslash escapes are skipped.
//
// For single quotes it distinguishes Rust char literals (e.g. 'a', '\n') from
// lifetimes (e.g. 'a in &'a str): a lifetime has no nearby closing quote, so
// start is returned unchanged and the quote is treated as an ordinary char.
// If a string literal is unterminated, the index of the last character is returned.
func skipQuoted(line string, start int) int {
	quote := line[start]
	for i := start + 1; i < len(line); i++ {
		if line[i] == '\\' {
			i++ // skip the escaped character
			continue
		}
		if line[i] == quote {
			return i
		}
		// A single quote with no close within a few chars is a lifetime.
		if quote == '\'' && i-start > 3 {
			return start
		}
	}
	if quote == '\'' {
		return start // unterminated single quote: treat as lifetime/ordinary
	}
	return len(line) - 1
}
