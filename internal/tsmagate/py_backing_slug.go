//ff:func feature=gate type=helper control=iteration dimension=1 lang=python
//ff:what pyBackingSlug: turns an item key (e.g. a qualified name "src.calc.classify") into a valid Python module identifier fragment by replacing every non [A-Za-z0-9_] rune with "_". Stricter than the TS slug because pytest imports the backing test BY MODULE NAME (prepend mode) — a dash or dot in the filename would make the module unimportable, so the backing filename must be a legal identifier.
package tsmagate

// pyBackingSlug replaces all non-identifier runes in key with underscores so the
// resulting gen_<slug>.py is an importable Python module name.
func pyBackingSlug(key string) string {
	out := []rune(key)
	for i, r := range out {
		isAlnum := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
		if !isAlnum {
			out[i] = '_'
		}
	}
	return string(out)
}
