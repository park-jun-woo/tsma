//ff:func feature=runner type=helper control=sequence
//ff:what Constructs command-line arguments for `mvn` to run a single JUnit test class
package runner

// buildMavenTestArgs constructs the arguments for `mvn -Dtest=<Class> test`,
// running only the named test class. The function is environment-independent so
// it can be unit tested without a Maven installation.
func buildMavenTestArgs(testClass string) []string {
	return []string{"-Dtest=" + testClass, "test"}
}
