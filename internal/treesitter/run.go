//ff:func feature=index type=helper control=iteration dimension=1
//ff:what Run: invokes `tree-sitter parse --xml` over a batch of files in one subprocess and returns the XML on stdout. Stderr (the "no parser directories" warning) is discarded; a non-zero exit (tree-sitter returns 1 when a file has ERROR nodes) is tolerated as long as XML was produced. A transient empty-output exit (cold grammar-cache compile race under concurrent processes) is retried a bounded number of times.
package treesitter

import "time"

// treeSitterRunAttempts bounds how many times Run re-invokes the CLI on a
// transient empty-output failure. treeSitterRetryBackoff is the inter-attempt
// pause; it is a var so tests can zero it for speed.
const treeSitterRunAttempts = 3

var treeSitterRetryBackoff = 20 * time.Millisecond

// Run executes the tree-sitter CLI to parse absPaths and returns the raw `--xml`
// output. grammarDir, when non-empty, is passed via -p so the grammar need not be
// globally configured. The captured stdout is returned even on a non-zero exit
// (parse errors); a transient empty-output process failure is retried, and only a
// persistent empty result surfaces as an error (callers then graceful-fallback).
func Run(command, grammarDir string, absPaths []string) ([]byte, error) {
	args := []string{"parse"}
	if grammarDir != "" {
		args = append(args, "-p", grammarDir)
	}
	args = append(args, "--xml")
	args = append(args, absPaths...)

	var out []byte
	var err error
	for attempt := 0; attempt < treeSitterRunAttempts; attempt++ {
		out, err = runOnce(command, args)
		if !shouldRetryRun(out, err) {
			break
		}
		time.Sleep(treeSitterRetryBackoff)
	}
	if len(out) == 0 && err != nil {
		return nil, err
	}
	return out, nil
}
