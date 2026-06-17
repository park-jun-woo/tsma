//ff:func feature=index type=helper control=sequence
//ff:what shouldRetryRun: decides whether a tree-sitter invocation failed transiently and is worth retrying. True only when the process actually ran and exited non-zero with NO stdout — the signature of the cold grammar-cache compile race. A start failure (CLI absent) is not an ExitError, so it fast-fails to caller fallback instead of retrying.
package treesitter

import (
	"os/exec"

	"errors"
)

// shouldRetryRun reports whether (out, err) is a transient empty-output process
// failure. A non-empty stdout (valid XML, possibly with ERROR nodes) or a nil
// error is success; a non-ExitError (e.g. command-not-found) is permanent.
func shouldRetryRun(out []byte, err error) bool {
	if err == nil || len(out) > 0 {
		return false
	}
	var exitErr *exec.ExitError
	return errors.As(err, &exitErr)
}
