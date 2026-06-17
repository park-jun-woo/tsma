//ff:func feature=match type=helper control=sequence lang=python
//ff:what canonicalPyTestPath: the Python arm of CanonicalTestPath (pytest) — foo.py → test_foo.py in the same directory (Phase005b §5 decision: same dir, not tests/); non-.py files return "". Same-dir + test_ prefix matches PyMatcher.Match's 1st search dir so the loop's write path and the matcher's read path agree on one convention.
package match

import (
	"path/filepath"
	"strings"
)

// canonicalPyTestPath returns test_<base> in the source's directory, or "" when
// base is not a .py source.
func canonicalPyTestPath(sourceFile, base string) string {
	if !strings.HasSuffix(base, ".py") {
		return ""
	}
	return filepath.Join(filepath.Dir(sourceFile), "test_"+base)
}
