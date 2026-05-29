//ff:func feature=runner type=helper control=sequence
//ff:what Constructs command-line arguments for `gradle` to run a single JUnit test class
package runner

// buildGradleTestArgs constructs the arguments for `gradle test --tests
// <Class>`, running only the named test class. The function is
// environment-independent so it can be unit tested without a Gradle
// installation.
func buildGradleTestArgs(testClass string) []string {
	return []string{"test", "--tests", testClass}
}
