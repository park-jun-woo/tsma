//ff:func feature=match type=helper control=sequence lang=python
//ff:what collectPyCalledNames: runs the ast names script against one test file and returns the referenced names as a slice. The Python analogue of collectTSCalledNames, but via the `ast` subprocess (parent §7-1: Python uses built-in ast, not tree-sitter). A subprocess or JSON error yields nil so one unparseable test file never aborts the package index.
package match

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

// collectPyCalledNames returns the names referenced (called or imported) by the
// test file at absPath, or nil on any error.
func collectPyCalledNames(python, absPath string) []string {
	cmd := exec.Command(python, "-c", pyAstNamesScript, absPath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil
	}
	var names []string
	if err := json.Unmarshal(stdout.Bytes(), &names); err != nil {
		return nil
	}
	return names
}
