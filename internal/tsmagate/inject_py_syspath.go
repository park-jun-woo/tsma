//ff:func feature=gate type=helper control=sequence lang=python
//ff:what injectPySyspath: prepends a sys.path bootstrap to a generated test so that, when the test runs from the isolation dir (.tsma/test) instead of beside the source, an `from <module> import foo` still resolves to the REAL module. It inserts the source directory at sys.path[0] at import time. This is the Python equivalent of rewriteTSImports — but Python imports are module names (not relative file paths), so injecting the source dir is the right seam (plan §5: "sys.path 주입/절대 import"). Only the backing copy gets this header; the canonical promote keeps the original (same-dir pytest resolves it natively).
package tsmagate

import "strconv"

// injectPySyspath returns src with a sys.path.insert(0, sourceDirAbs) header so a
// backing test under .tsma/test imports the real source module by name.
func injectPySyspath(src, sourceDirAbs string) string {
	header := "import sys as _tsma_sys\n_tsma_sys.path.insert(0, " + strconv.Quote(sourceDirAbs) + ")\n"
	return header + src
}
