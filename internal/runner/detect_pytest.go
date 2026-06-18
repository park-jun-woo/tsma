//ff:func feature=runner type=helper control=sequence
//ff:what Delegates to detect.DetectPytest (SSOT) to decide if pytest is used
package runner

import "github.com/park-jun-woo/tsma/internal/detect"

// detectPytest reports whether the project uses pytest. It is a thin wrapper
// over detect.DetectPytest, the single source of truth shared with the coverage
// stage (see BUG-001 / Phase006: the runner used to fall back to unittest while
// coverage already assumed pytest). The package-local containsPytest/fileExists
// helpers are retained for their unit tests; the detection body now lives in the
// detect package.
func detectPytest(projectRoot string) bool {
	return detect.DetectPytest(projectRoot)
}
