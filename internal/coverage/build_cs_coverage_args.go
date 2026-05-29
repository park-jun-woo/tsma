//ff:func feature=coverage type=helper control=sequence lang=csharp
//ff:what Builds dotnet test arguments that collect Cobertura coverage into the given results directory
package coverage

// buildCsCoverageArgs returns the `dotnet test` arguments that run the test
// suite and collect coverage via coverlet's XPlat collector, emitting a
// Cobertura report under resultsDir. The collector writes the report to
// <resultsDir>/<guid>/coverage.cobertura.xml. Kept environment-independent for
// unit testing.
func buildCsCoverageArgs(resultsDir string) []string {
	return []string{
		"test",
		"--collect:XPlat Code Coverage",
		"--results-directory", resultsDir,
	}
}
