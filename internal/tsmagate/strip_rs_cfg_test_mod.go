//ff:func feature=gate type=helper control=sequence lang=rust
//ff:what stripRsCfgTestMod: removes a trailing `#[cfg(test)] mod ... { ... }` block from a Rust source so a freshly generated test module replaces (not duplicates — a duplicate `mod tests` is a compile error) the previous one. Rust convention places the unit-test module at the end of the file, so cutting from the last `#[cfg(test)]` marker to EOF is the precise, regenerate-safe seam; a source with no in-file test module is returned unchanged (brand-new function → the generated block is simply appended by injectRsTestMod).
package tsmagate

import "strings"

// stripRsCfgTestMod returns src with any trailing #[cfg(test)] module removed
// (back to the last #[cfg(test)] marker), or src unchanged when none is present.
func stripRsCfgTestMod(src string) string {
	idx := strings.LastIndex(src, "#[cfg(test)]")
	if idx < 0 {
		return src
	}
	return strings.TrimRight(src[:idx], " \t\n") + "\n"
}
