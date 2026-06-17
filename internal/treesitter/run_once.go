//ff:func feature=index type=helper control=sequence
//ff:what runOnce: executes the tree-sitter CLI once over args, capturing stdout and discarding the "no parser directories" stderr warning. Returns captured stdout (even on non-zero exit) and the process error; the retry policy lives in Run.
package treesitter

import (
	"bytes"
	"io"
	"os/exec"
)

// runOnce performs a single `tree-sitter` invocation and returns its stdout
// bytes alongside the raw command error (nil on exit 0).
func runOnce(command string, args []string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = io.Discard
	err := cmd.Run()
	return stdout.Bytes(), err
}
