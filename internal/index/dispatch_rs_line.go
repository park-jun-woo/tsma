//ff:func feature=index type=helper control=selection
//ff:what Dispatches a Rust source line to impl/mod scope push or fn declaration handling
package index

// dispatchRsLine classifies the current line as an impl block, module, or
// function declaration and updates the parse state accordingly.
func dispatchRsLine(st *rsParseState, trimmed, line string) {
	switch {
	case isRsImplLine(trimmed, line):
		m := rsImplPattern.FindStringSubmatch(trimmed)
		closePrevTSEndLine(st.functions, st.lineNum, st.lastNonEmptyLine)
		st.scopes = append(st.scopes, rsScope{depth: st.depth, receiver: m[1], cfgTest: cfgTestActive(st.scopes, st.pendingCfgTest)})
	case rsModPattern.MatchString(trimmed):
		m := rsModPattern.FindStringSubmatch(trimmed)
		st.scopes = append(st.scopes, rsScope{depth: st.depth, module: m[1], cfgTest: cfgTestActive(st.scopes, st.pendingCfgTest)})
	case rsFnPattern.MatchString(trimmed):
		appendRsFunc(st, trimmed)
	}
}
