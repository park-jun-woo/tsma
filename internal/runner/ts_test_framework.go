//ff:type feature=runner type=model
//ff:what Represents a detected TypeScript test framework (vitest, jest, mocha)
package runner

// tsTestFramework represents a detected TypeScript test framework.
type tsTestFramework string

const (
	frameworkVitest tsTestFramework = "vitest"
	frameworkJest   tsTestFramework = "jest"
	frameworkMocha  tsTestFramework = "mocha"
)
