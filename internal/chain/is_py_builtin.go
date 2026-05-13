//ff:func feature=chain type=helper control=sequence
//ff:what Returns true for common Python builtins that should be skipped during tracing
package chain

// isPyBuiltin returns true for common Python builtins that should be skipped.
func isPyBuiltin(name string) bool {
	builtins := map[string]bool{
		"print": true, "len": true, "range": true, "str": true,
		"int": true, "float": true, "bool": true, "list": true,
		"dict": true, "set": true, "tuple": true, "type": true,
		"isinstance": true, "issubclass": true, "hasattr": true,
		"getattr": true, "setattr": true, "delattr": true,
		"super": true, "property": true, "staticmethod": true,
		"classmethod": true, "enumerate": true, "zip": true,
		"map": true, "filter": true, "sorted": true, "reversed": true,
		"any": true, "all": true, "min": true, "max": true,
		"sum": true, "abs": true, "round": true, "repr": true,
		"format": true, "open": true, "input": true,
		"raise": true, "return": true, "yield": true,
	}
	return builtins[name]
}
