//ff:func feature=index type=helper control=sequence
//ff:what Records pending #[cfg(test)] attributes and reports whether the line was an attribute
package index

import "strings"

// handleRsAttribute inspects a trimmed line for an attribute (`#[...]`). If it
// is a #[cfg(test)] attribute it marks the parse state so the guarded item is
// treated as test-only. It returns true when the line was an attribute and the
// caller should skip further processing.
func handleRsAttribute(st *rsParseState, trimmed string) bool {
	if !strings.HasPrefix(trimmed, "#[") {
		return false
	}
	if strings.Contains(trimmed, "cfg(test)") {
		st.pendingCfgTest = true
	}
	return true
}
