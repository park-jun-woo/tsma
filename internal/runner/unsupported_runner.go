//ff:type feature=runner type=model
//ff:what Placeholder runner for languages without test runner support
package runner

// UnsupportedRunner is a placeholder for unimplemented languages.
type UnsupportedRunner struct {
	Lang string
}
