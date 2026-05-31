//ff:test feature=cli
package cli

import "testing"

func TestFirstLine(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"single line no newline", "hello", "hello"},
		{"multi line", "first\nsecond\nthird", "first"},
		{"leading newline", "\nrest", ""},
		{"trailing newline only", "only\n", "only"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firstLine(tt.in); got != tt.want {
				t.Errorf("firstLine(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
