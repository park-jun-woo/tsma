//ff:func feature=graph type=helper control=selection
//ff:what Returns true for common JS/TS built-in objects that should be skipped
package graph

// isTSBuiltin returns true for common JS/TS built-in objects.
func isTSBuiltin(name string) bool {
	switch name {
	case "console", "JSON", "Math", "Date", "Array", "Object", "String",
		"Number", "Boolean", "RegExp", "Error", "Promise", "Map", "Set",
		"parseInt", "parseFloat", "setTimeout", "setInterval", "clearTimeout",
		"clearInterval", "require", "process":
		return true
	}
	return false
}
