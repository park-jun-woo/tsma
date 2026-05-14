package match

import "testing"

func TestIsTSTestFileVariants(t *testing.T) {
	tests := []struct {
		name string
		file string
		want bool
	}{
		{name: "test.ts", file: "handler.test.ts", want: true},
		{name: "test.js", file: "handler.test.js", want: true},
		{name: "spec.ts", file: "handler.spec.ts", want: true},
		{name: "spec.js", file: "handler.spec.js", want: true},
		{name: "regular ts", file: "handler.ts", want: false},
		{name: "regular js", file: "handler.js", want: false},
		{name: "go file", file: "handler_test.go", want: false},
		{name: "py file", file: "test_handler.py", want: false},
		{name: "empty", file: "", want: false},
		{name: "just suffix", file: ".test.ts", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTSTestFile(tt.file)
			if got != tt.want {
				t.Errorf("isTSTestFile(%q) = %v, want %v", tt.file, got, tt.want)
			}
		})
	}
}
