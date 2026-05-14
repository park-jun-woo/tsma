package coverage

import "testing"

func TestExtractCoverFile(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "valid position part",
			input: "github.com/example/pkg/handler.go:10.2,20.5",
			want:  "github.com/example/pkg/handler.go",
		},
		{
			name:  "nested colons in path",
			input: "C:/Users/example/handler.go:10.2,20.5",
			want:  "C:/Users/example/handler.go",
		},
		{
			name:    "no colon",
			input:   "nocolonhere",
			wantErr: true,
		},
		{
			name:    "no comma after colon",
			input:   "handler.go:10.2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractCoverFile(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
