package match

import "testing"

func TestUnsupportedMatcherAlwaysReturnsFalse(t *testing.T) {
	m := &unsupportedMatcher{}

	tests := []struct {
		name       string
		root       string
		sourceFile string
	}{
		{name: "rust file", root: "/tmp", sourceFile: "handler.rs"},
		{name: "java file", root: "/project", sourceFile: "Main.java"},
		{name: "empty", root: "", sourceFile: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile, found := m.Match(tt.root, tt.sourceFile)
			if found {
				t.Errorf("expected found=false for %q", tt.sourceFile)
			}
			if testFile != "" {
				t.Errorf("expected empty testFile, got %q", testFile)
			}
		})
	}
}
