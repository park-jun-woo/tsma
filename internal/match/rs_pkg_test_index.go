//ff:type feature=match type=model lang=rust
//ff:what RsTestIndex maps bare source function names to the test files that call them — the content-aware index for one Rust source file's test candidates. Keys come from the file's own in-file #[cfg(test)] module (attributing to the source file itself) and from any tests/<name>.rs integration test (attributing to that file). Coarser than the Go index (file granularity, no receiver): cargo runs whole test binaries, so TestFuncs is left nil downstream.
package match

// RsTestIndex is a content-aware index for a single Rust source file's test
// candidates. Keys are bare function names called by the in-file #[cfg(test)]
// module or by tests/<name>.rs; each maps to the test files that reference it.
type RsTestIndex struct {
	refs map[string][]string
}
