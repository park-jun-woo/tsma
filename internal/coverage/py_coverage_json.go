//ff:type feature=coverage type=model
//ff:what Represents the top-level structure of coverage.py JSON output
package coverage

// pyCoverageJSON represents the top-level structure of coverage.py JSON output.
type pyCoverageJSON struct {
	Files map[string]pyCoverageFile `json:"files"`
}
