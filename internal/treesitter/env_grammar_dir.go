//ff:func feature=index type=helper control=sequence
//ff:what envGrammarDir: validates an explicit TSMA_<LANG>_GRAMMAR override — returns it when it is an existing directory, else "" (a bad override falls through to line-based rather than erroring).
package treesitter

// envGrammarDir returns g when it is an existing directory, else "".
func envGrammarDir(g string) string {
	if grammarDirExists(g) {
		return g
	}
	return ""
}
