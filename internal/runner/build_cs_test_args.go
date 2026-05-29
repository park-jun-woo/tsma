//ff:func feature=runner type=helper control=sequence lang=csharp
//ff:what Builds dotnet test command arguments that filter to a single C# test class
package runner

import "fmt"

// buildCsTestArgs constructs the arguments for `dotnet test` to run only the
// given test class, using the fully-qualified-name substring filter
// (FullyQualifiedName~Class). Kept environment-independent for unit testing.
func buildCsTestArgs(testClass string) []string {
	return []string{"test", "--filter", fmt.Sprintf("FullyQualifiedName~%s", testClass)}
}
