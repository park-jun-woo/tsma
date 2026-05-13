//ff:type feature=coverage type=model
//ff:what Describes a branch that was not covered by tests
package coverage

// UncoveredBranch describes a branch that was not covered.
type UncoveredBranch struct {
	File string
	Line int
}
