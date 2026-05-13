//ff:func feature=coverage type=helper control=selection
//ff:what Builds npx argument list for the detected framework with JSON coverage output
package coverage

// buildCoverageArgs constructs the command-line arguments for running tests with coverage.
func buildCoverageArgs(framework tsCoverageFramework, testFile, coverDir string) []string {
	switch framework {
	case coverJest:
		return []string{
			"jest", testFile,
			"--coverage",
			"--coverageReporters=json",
			"--coverageDirectory=" + coverDir,
		}
	case coverVitest:
		return []string{
			"vitest", "run", testFile,
			"--coverage",
			"--coverage.provider=v8",
			"--coverage.reporter=json",
			"--coverage.reportsDirectory=" + coverDir,
		}
	default:
		return []string{
			"vitest", "run", testFile,
			"--coverage",
			"--coverage.provider=v8",
			"--coverage.reporter=json",
			"--coverage.reportsDirectory=" + coverDir,
		}
	}
}
