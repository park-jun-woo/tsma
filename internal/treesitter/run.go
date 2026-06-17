//ff:func feature=index type=helper control=sequence
//ff:what Run: invokes `tree-sitter parse --xml` over a batch of files in one subprocess and returns the XML on stdout. Stderr (the "no parser directories" warning) is discarded; a non-zero exit (tree-sitter returns 1 when a file has ERROR nodes) is tolerated as long as XML was produced.
package treesitter

import (
	"bytes"
	"io"
	"os/exec"
)

// Run executes the tree-sitter CLI to parse absPaths and returns the raw `--xml`
// output. grammarDir, when non-empty, is passed via -p so the grammar need not be
// globally configured. The captured stdout is returned even on a non-zero exit
// (parse errors), erroring only when nothing was produced.
func Run(command, grammarDir string, absPaths []string) ([]byte, error) {
	args := []string{"parse"}
	if grammarDir != "" {
		args = append(args, "-p", grammarDir)
	}
	args = append(args, "--xml")
	args = append(args, absPaths...)

	cmd := exec.Command(command, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = io.Discard

	err := cmd.Run()
	out := stdout.Bytes()
	if len(out) == 0 && err != nil {
		return nil, err
	}
	return out, nil
}
