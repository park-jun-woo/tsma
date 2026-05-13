//ff:func feature=runner type=helper control=sequence
//ff:what Constructs command-line arguments for go test
package runner

import "strings"

// buildGoTestArgs constructs command-line arguments for go test.
func buildGoTestArgs(pkgPath string, testFuncs []string) []string {
	args := []string{"test", "-v", "-count=1"}
	if len(testFuncs) > 0 {
		args = append(args, "-run", strings.Join(testFuncs, "|"))
	}
	args = append(args, pkgPath)
	return args
}
