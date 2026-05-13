//ff:func feature=chain type=helper control=sequence
//ff:what Constructs the index key for a Python function based on module, class, and indent
package chain

// buildPyFuncKey constructs the index key for a Python function.
func buildPyFuncKey(modKey, funcName string, indentLen int, currentClass string, classIndentLen int) string {
	if currentClass != "" && indentLen > classIndentLen {
		return modKey + "." + currentClass + "." + funcName
	}
	return modKey + "." + funcName
}
