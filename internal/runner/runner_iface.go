//ff:type feature=runner type=model
//ff:what Defines the Runner interface for language-specific test execution
package runner

// Runner executes tests for a given language.
type Runner interface {
	Run(projectRoot, testFile string) (*Result, error)
}
