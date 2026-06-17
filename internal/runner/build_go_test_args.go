//ff:func feature=runner type=helper control=sequence
//ff:what Constructs command-line arguments for go test
package runner

// buildGoTestArgs constructs command-line arguments for go test. When test
// functions are given, the -run pattern is anchored per name (^Name$) and
// alternated so that e.g. TestGenerate does not also match TestGenerateBytes.
// overlayArgs (e.g. -overlay/-vet=off from GoOverlayArgs) are inserted before the
// package path; nil for the manual-submit/disk-truth path.
func buildGoTestArgs(pkgPath string, testFuncs []string, overlayArgs []string) []string {
	args := []string{"test", "-v", "-count=1"}
	if pattern := AnchorRunPattern(testFuncs); pattern != "" {
		args = append(args, "-run", pattern)
	}
	args = append(args, overlayArgs...)
	args = append(args, pkgPath)
	return args
}
