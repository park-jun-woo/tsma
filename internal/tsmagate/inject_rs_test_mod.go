//ff:func feature=gate type=helper control=sequence lang=rust
//ff:what injectRsTestMod: produces the new source-file contents for the D5 in-file mod strategy — strips any previous trailing #[cfg(test)] module (stripRsCfgTestMod, regenerate-safe) and appends the generated test block, wrapping it in `#[cfg(test)] mod tests { ... }` only when the LLM did not already emit the cfg(test) guard. The result is written to the SOURCE file during measurement (and rolled back on a non-terminal result), so a non-pub function's in-file unit test can exercise it — the capability the tests/ integration path lacks.
package tsmagate

import "strings"

// injectRsTestMod returns originalSrc with its trailing #[cfg(test)] module
// replaced by (or, for a brand-new function, the appended) generated test block.
func injectRsTestMod(originalSrc, genBlock string) string {
	base := stripRsCfgTestMod(originalSrc)
	block := strings.TrimSpace(genBlock)
	if !strings.Contains(block, "#[cfg(test)]") {
		block = "#[cfg(test)]\nmod tests {\n" + block + "\n}"
	}
	return strings.TrimRight(base, "\n") + "\n\n" + block + "\n"
}
