//ff:func feature=chain type=helper control=sequence
//ff:what Extracts displayName, funcName, and receiver presence from TS call match groups
package chain

// parseTSCallNames extracts displayName, funcName, and receiver presence from match groups.
func parseTSCallNames(receiver, method, plainFunc string) (string, string, bool) {
	if receiver != "" && method != "" {
		if tsSkipCallNames[receiver] {
			return "", "", false
		}
		return receiver + "." + method, method, true
	}
	if plainFunc != "" {
		if tsSkipCallNames[plainFunc] {
			return "", "", false
		}
		return plainFunc, plainFunc, false
	}
	return "", "", false
}
