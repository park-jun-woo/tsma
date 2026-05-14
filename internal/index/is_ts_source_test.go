package index

import "testing"

func TestIsTSSourceVariants(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"handler.ts", true},
		{"handler.js", true},
		{"component.tsx", true},
		{"page.jsx", true},
		{"component.test.tsx", false},
		{"component.spec.jsx", false},
		{"types.d.ts", false},
		{"handler.test.ts", false},
		{"handler.spec.ts", false},
		{"handler.test.js", false},
		{"handler.spec.js", false},
		{"readme.md", false},
		{"handler.go", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isTSSource(tt.path)
			if got != tt.want {
				t.Errorf("isTSSource(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
