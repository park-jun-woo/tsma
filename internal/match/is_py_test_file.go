//ff:func feature=match type=helper control=sequence lang=python
//ff:what isPyTestFile: recognizes a pytest/unittest test module by the two conventional shapes — test_*.py or *_test.py. Used by ingestPyDir to pick which files feed the content-aware index. Non-.py names are rejected.
package match

import "strings"

// isPyTestFile returns true for the pytest naming conventions test_<x>.py and
// <x>_test.py.
func isPyTestFile(name string) bool {
	if !strings.HasSuffix(name, ".py") {
		return false
	}
	return strings.HasPrefix(name, "test_") || strings.HasSuffix(name, "_test.py")
}
