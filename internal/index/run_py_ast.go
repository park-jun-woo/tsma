//ff:func feature=index type=helper control=sequence lang=python
//ff:what runPyAst: runs the embedded ast script (`python -c <script> <file>`) and returns its stdout JSON. The single subprocess seam for Python D1 — the analogue of treesitter.Run. A non-zero exit (SyntaxError → exit 2, or any interpreter error) is returned as an error so the caller per-file-falls-back to the line indexer; stderr is folded into the error for diagnosis.
package index

import (
	"bytes"
	"fmt"
	"os/exec"
)

// runPyAst executes the ast dump script against absPath and returns stdout. An
// interpreter or parse failure is returned as an error (caller falls back).
func runPyAst(python, script, absPath string) ([]byte, error) {
	cmd := exec.Command(python, "-c", script, absPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("python ast parse %s: %v: %s", absPath, err, stderr.String())
	}
	return stdout.Bytes(), nil
}
