//ff:func feature=match type=helper control=sequence lang=rust
//ff:what canonicalRsTestPath: the Rust arm of CanonicalTestPath. Rust unit tests live INSIDE the source file in a `#[cfg(test)] mod tests` block, not in a separate file — so the canonical "test path" for a .rs source is the source file itself (the loop's D5 in-file mod injection writes there, and RsMatcher already attributes an in-file test back to the source file, so write path and read path agree). A non-.rs file returns "".
package match

import "strings"

// canonicalRsTestPath returns the Rust source file as its own test path (in-file
// #[cfg(test)] mod), or "" when base is not a .rs file.
func canonicalRsTestPath(sourceFile, base string) string {
	if !strings.HasSuffix(base, ".rs") {
		return ""
	}
	return sourceFile
}
