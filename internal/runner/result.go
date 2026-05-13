//ff:type feature=runner type=model
//ff:what Holds the outcome of a test run including pass/fail status and output
package runner

// Result holds the outcome of a test run.
type Result struct {
	Pass   bool
	Output string
}
