//ff:func feature=match type=helper control=sequence lang=python
//ff:what resolvePython: locates the Python interpreter for the content-aware matcher's ast subprocess, preferring python3 over python. Returns "" when neither is on PATH — the signal PythonFuncMatcher uses to fall back to filename matching (graceful). A package-local copy (mirrors runner/coverage's own findPython) so match never imports those packages.
package match

import "os/exec"

// resolvePython returns "python3", then "python", then "" when no interpreter
// is on PATH (which triggers the filename fallback).
func resolvePython() string {
	if _, err := exec.LookPath("python3"); err == nil {
		return "python3"
	}
	if _, err := exec.LookPath("python"); err == nil {
		return "python"
	}
	return ""
}
