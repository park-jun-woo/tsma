package index

import "testing"

func TestIsGoSourceVariants(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"handler.go", true},
		{"internal/api/handler.go", true},
		{"handler_test.go", false},
		{"internal/api/handler_test.go", false},
		{"mock_service.go", false},
		{"internal/mock_repo.go", false},
		{"server_gen.go", false},
		{"openapi.gen.go", false},
		{"message.pb.go", false},
		{"internal/api/server_gen.go", false},
		{"internal/api/types.gen.go", false},
		{"internal/proto/msg.pb.go", false},
		{"readme.md", false},
		{"main.py", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isGoSource(tt.path)
			if got != tt.want {
				t.Errorf("isGoSource(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
