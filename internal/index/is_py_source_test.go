package index

import "testing"

func TestIsPySourceVariants(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"handler.py", true},
		{"app/auth.py", true},
		{"test_handler.py", false},
		{"app/test_auth.py", false},
		{"readme.md", false},
		{"handler.go", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isPySource(tt.path)
			if got != tt.want {
				t.Errorf("isPySource(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
