package index

import "testing"

func TestIsCsControlKeyword(t *testing.T) {
	for _, kw := range []string{"if", "for", "foreach", "while", "switch", "using", "lock", "get", "set"} {
		if !isCsControlKeyword(kw) {
			t.Errorf("isCsControlKeyword(%q) = false, want true", kw)
		}
	}
	for _, name := range []string{"Add", "Classify", "Foo", "Compute"} {
		if isCsControlKeyword(name) {
			t.Errorf("isCsControlKeyword(%q) = true, want false", name)
		}
	}
}
