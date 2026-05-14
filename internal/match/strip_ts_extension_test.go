package match

import "testing"

func TestStripTSExtensionVariants(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: ".ts", input: "handler.ts", want: "handler"},
		{name: ".js", input: "handler.js", want: "handler"},
		{name: ".tsx", input: "Component.tsx", want: "Component"},
		{name: ".jsx", input: "Component.jsx", want: "Component"},
		{name: "no extension", input: "handler", want: "handler"},
		{name: "unknown ext", input: "handler.py", want: "handler.py"},
		{name: "multiple dots ts", input: "my.handler.ts", want: "my.handler"},
		{name: "multiple dots tsx", input: "my.component.tsx", want: "my.component"},
		{name: "empty string", input: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripTSExtension(tt.input)
			if got != tt.want {
				t.Errorf("stripTSExtension(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
