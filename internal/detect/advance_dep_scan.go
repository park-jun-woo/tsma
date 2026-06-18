//ff:func feature=detect type=helper control=sequence lang=python
//ff:what Advances the PEP 621 dep-scan state for one line; reports a pytest hit
package detect

import "strings"

// advanceDepScan consumes one already-trimmed pyproject.toml line and returns
// whether it declares pytest (hit) plus the updated scan state. It is a pure
// transition function so containsPytestDep stays a flat single-loop scan (Q1).
//
//   - a [section] header re-evaluates the dependency-table flag and closes any
//     open dependency array;
//   - a `dependencies = [...]` assignment reports immediately on an inline pytest
//     token, else opens a multiline array;
//   - inside a dependency table or array, a pytest token is a hit; a closing `]`
//     ends a multiline array.
func advanceDepScan(st depScanState, line string) (bool, depScanState) {
	lower := strings.ToLower(line)

	if strings.HasPrefix(line, "[") {
		st.inDepArray = false
		st.inDepTable = isDepTableHeader(line)
		return false, st
	}

	if strings.HasPrefix(lower, "dependencies") && strings.Contains(line, "[") {
		if strings.Contains(lower, "pytest") {
			return true, st
		}
		st.inDepArray = !strings.Contains(line, "]")
		return false, st
	}

	if (st.inDepTable || st.inDepArray) && strings.Contains(lower, "pytest") {
		return true, st
	}
	if st.inDepArray && strings.Contains(line, "]") {
		st.inDepArray = false
	}
	return false, st
}
