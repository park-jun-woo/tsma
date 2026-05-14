package index

import "testing"

func TestResetTSClassContext(t *testing.T) {
	tests := []struct {
		name        string
		trimmed     string
		line        string
		classIndent int
		want        bool
	}{
		{
			name:        "closing brace at class level",
			trimmed:     "}",
			line:        "}",
			classIndent: 0,
			want:        true,
		},
		{
			name:        "empty line does not reset",
			trimmed:     "",
			line:        "",
			classIndent: 0,
			want:        false,
		},
		{
			name:        "indented line inside class does not reset",
			trimmed:     "login() {",
			line:        "  login() {",
			classIndent: 0,
			want:        false,
		},
		{
			name:        "new class declaration does not reset",
			trimmed:     "class AnotherClass {",
			line:        "class AnotherClass {",
			classIndent: 0,
			want:        false,
		},
		{
			name:        "non-class line at class indent resets (not closing brace)",
			trimmed:     "function standalone() {",
			line:        "function standalone() {",
			classIndent: 0,
			want:        false,
		},
		{
			name:        "closing brace deeper than class does not reset",
			trimmed:     "}",
			line:        "    }",
			classIndent: 0,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resetTSClassContext(tt.trimmed, tt.line, tt.classIndent)
			if got != tt.want {
				t.Errorf("resetTSClassContext(%q, %q, %d) = %v, want %v",
					tt.trimmed, tt.line, tt.classIndent, got, tt.want)
			}
		})
	}
}
