//ff:type feature=match type=factory
//ff:what Defines the Matcher interface for source-file-to-test-file matching
package match

// Matcher finds the test file that covers a given source file.
type Matcher interface {
	Match(projectRoot string, sourceFile string) (testFile string, found bool)
}
