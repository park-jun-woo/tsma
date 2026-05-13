//ff:type feature=coverage type=model
//ff:what Represents a detected framework for TS/JS coverage purposes (vitest or jest)
package coverage

// tsCoverageFramework represents a detected framework for coverage purposes.
type tsCoverageFramework string

const (
	coverVitest tsCoverageFramework = "vitest"
	coverJest   tsCoverageFramework = "jest"
)
