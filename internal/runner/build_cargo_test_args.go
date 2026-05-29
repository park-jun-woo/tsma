//ff:func feature=runner type=helper control=sequence
//ff:what Constructs command-line arguments for cargo test from a test file path
package runner

import (
	"path/filepath"
	"strings"
)

// buildCargoTestArgs constructs the arguments for `cargo test` given the
// project-relative test file.
//
//   - Integration tests live in tests/<name>.rs and run via `--test <name>`.
//   - In-file unit tests (any other .rs file) run with a plain `cargo test`,
//     which executes the crate's unit tests.
//
// The function is environment-independent so it can be unit tested without cargo.
func buildCargoTestArgs(testFile string) []string {
	args := []string{"test"}

	dir := filepath.ToSlash(filepath.Dir(testFile))
	if dir == "tests" || strings.HasPrefix(dir, "tests/") {
		name := strings.TrimSuffix(filepath.Base(testFile), ".rs")
		args = append(args, "--test", name)
	}

	return args
}
