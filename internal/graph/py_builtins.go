//ff:func feature=graph type=helper control=sequence
//ff:what Defines the set of Python built-in functions to skip during analysis
package graph

// pyBuiltins are Python built-in functions to skip.
var pyBuiltins = map[string]bool{
	"print": true, "len": true, "range": true, "str": true, "int": true,
	"float": true, "bool": true, "list": true, "dict": true, "set": true,
	"tuple": true, "type": true, "isinstance": true, "issubclass": true,
	"super": true, "hasattr": true, "getattr": true, "setattr": true,
	"enumerate": true, "zip": true, "map": true, "filter": true,
	"sorted": true, "reversed": true, "min": true, "max": true, "sum": true,
	"any": true, "all": true, "abs": true, "round": true, "open": true,
	"input": true, "repr": true, "id": true, "hash": true, "hex": true,
	"oct": true, "bin": true, "chr": true, "ord": true, "next": true,
	"iter": true, "format": true, "vars": true, "dir": true, "help": true,
}
