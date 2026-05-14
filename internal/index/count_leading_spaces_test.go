package index

import "testing"

func TestCountLeadingSpaces(t *testing.T) {
	tests := []struct {
		name string
		line string
		want int
	}{
		{name: "no indent", line: "hello", want: 0},
		{name: "two spaces", line: "  hello", want: 2},
		{name: "four spaces", line: "    hello", want: 4},
		{name: "single tab", line: "\thello", want: 4},
		{name: "two tabs", line: "\t\thello", want: 8},
		{name: "mixed space and tab", line: "  \thello", want: 6},
		{name: "empty string", line: "", want: 0},
		{name: "only spaces", line: "    ", want: 4},
		{name: "only tabs", line: "\t\t", want: 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countLeadingSpaces(tt.line)
			if got != tt.want {
				t.Errorf("countLeadingSpaces(%q) = %d, want %d", tt.line, got, tt.want)
			}
		})
	}
}
