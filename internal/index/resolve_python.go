//ff:func feature=index type=helper control=sequence lang=python
//ff:what resolvePython: locates the Python interpreter for the ast precise path, preferring python3 over python. Returns "" when neither is on PATH — the signal PyAstIndexer uses to fall back to the line-based indexer (graceful, zero regression per Phase005b §1). Distinct from runner/coverage findPython (which default to "python") because D1 must be able to detect absence.
package index

import "os/exec"

// resolvePython returns "python3" if available, else "python", else "" when no
// interpreter is on PATH (which triggers the line-based fallback).
func resolvePython() string {
	if _, err := exec.LookPath("python3"); err == nil {
		return "python3"
	}
	if _, err := exec.LookPath("python"); err == nil {
		return "python"
	}
	return ""
}
