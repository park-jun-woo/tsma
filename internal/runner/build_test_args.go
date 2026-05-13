//ff:func feature=runner type=helper control=selection
//ff:what Builds npx argument list for the detected test framework with verbose output
package runner

// buildTestArgs constructs the command-line arguments for the detected test framework.
func buildTestArgs(framework tsTestFramework, testFile string) []string {
	switch framework {
	case frameworkJest:
		return []string{"jest", testFile, "--verbose"}
	case frameworkMocha:
		return []string{"mocha", testFile}
	case frameworkVitest:
		return []string{"vitest", "run", testFile, "--reporter=verbose"}
	default:
		return []string{"vitest", "run", testFile, "--reporter=verbose"}
	}
}
