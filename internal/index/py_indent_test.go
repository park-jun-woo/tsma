package index

import "testing"

func TestPyIndent(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{name: "no indent", input: "hello", want: 0},
		{name: "four spaces", input: "    def foo():", want: 4},
		{name: "eight spaces", input: "        pass", want: 8},
		{name: "single tab", input: "\tdef foo():", want: 4},
		{name: "two tabs", input: "\t\tpass", want: 8},
		{name: "mixed spaces and tab", input: "  \thello", want: 6},
		{name: "empty string", input: "", want: 0},
		{name: "only whitespace", input: "   ", want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pyIndent(tt.input)
			if got != tt.want {
				t.Errorf("pyIndent(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
