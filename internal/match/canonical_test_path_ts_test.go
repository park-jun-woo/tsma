package match

import "testing"

// TestCanonicalTestPathTypescript covers the Phase005a TS branch: foo.ts →
// foo.test.ts in the same directory, extension preserved, non-TS files "".
func TestCanonicalTestPathTypescript(t *testing.T) {
	cases := []struct {
		src  string
		want string
	}{
		{"src/math.ts", "src/math.test.ts"},
		{"a/b/widget.tsx", "a/b/widget.test.tsx"},
		{"lib/util.js", "lib/util.test.js"},
		{"comp.jsx", "comp.test.jsx"},
		{"src/math.go", ""}, // wrong extension for typescript
		{"src/README.md", ""},
	}
	for _, c := range cases {
		if got := CanonicalTestPath("typescript", c.src); got != c.want {
			t.Errorf("CanonicalTestPath(typescript, %q) = %q, want %q", c.src, got, c.want)
		}
	}
	// Rust (Phase005e): the source file itself is the in-file #[cfg(test)] target.
	if got := CanonicalTestPath("rust", "src/x.rs"); got != "src/x.rs" {
		t.Errorf("rust should map to the source file, got %q", got)
	}
	// A genuinely non-handled language still returns "".
	if got := CanonicalTestPath("ruby", "x.rb"); got != "" {
		t.Errorf("ruby should be unhandled, got %q", got)
	}
}
