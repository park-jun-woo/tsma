//ff:func feature=gate type=helper control=iteration dimension=1 lang=rust
//ff:what parseRsTestFuncs: the Rust truncation check + test-name extraction (the go/parser-based parseTestFuncs analogue, but Rust has no in-process parser so it is lexical). ok=false signals a truncated generation — an unbalanced `{`/`}` count means the model's `#[cfg(test)] mod tests { ... }` block was cut off, which the loop turns into the C3 truncated feedback instead of measuring. On balance it returns every `#[test] fn <name>` (including async/tokio test attributes) so the caller can reject a block that declares no test (nothing would run). Unlike Go there is no malformed-name rule — Rust test names are arbitrary snake_case.
package tsmagate

import (
	"regexp"
	"strings"
)

// rsTestFnPattern captures the function name following each #[test]-style
// attribute (plain #[test], #[tokio::test], #[async_std::test], ...).
var rsTestFnPattern = regexp.MustCompile(`#\[(?:[\w:]+::)?test\b[^\]]*\][\s\S]*?fn\s+([A-Za-z_]\w*)`)

// parseRsTestFuncs returns the names of the #[test] functions in a generated
// block and ok=false when the brace count is unbalanced (a truncated output).
func parseRsTestFuncs(src string) ([]string, bool) {
	if strings.Count(src, "{") != strings.Count(src, "}") {
		return nil, false
	}
	var names []string
	for _, m := range rsTestFnPattern.FindAllStringSubmatch(src, -1) {
		names = append(names, m[1])
	}
	return names, true
}
