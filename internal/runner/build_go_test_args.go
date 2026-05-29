//ff:func feature=runner type=helper control=sequence
//ff:what Constructs command-line arguments for go test
package runner

// buildGoTestArgs constructs command-line arguments for go test. When test
// functions are given, the -run pattern is anchored per name (^Name$) and
// alternated so that e.g. TestGenerate does not also match TestGenerateBytes.
func buildGoTestArgs(pkgPath string, testFuncs []string) []string {
	args := []string{"test", "-v", "-count=1"}
	if pattern := AnchorRunPattern(testFuncs); pattern != "" {
		args = append(args, "-run", pattern)
	}
	args = append(args, pkgPath)
	return args
}
