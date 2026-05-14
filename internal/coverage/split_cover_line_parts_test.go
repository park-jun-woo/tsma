package coverage

import "testing"

func TestSplitCoverLineParts(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		wantFile   string
		wantStmts  string
		wantCount  string
		wantErr    bool
	}{
		{
			name:      "valid line",
			line:      "github.com/example/pkg/handler.go:10.2,20.5 3 1",
			wantFile:  "github.com/example/pkg/handler.go:10.2,20.5",
			wantStmts: "3",
			wantCount: "1",
		},
		{
			name:      "zero count",
			line:      "pkg/main.go:1.1,5.5 2 0",
			wantFile:  "pkg/main.go:1.1,5.5",
			wantStmts: "2",
			wantCount: "0",
		},
		{
			name:    "no spaces",
			line:    "nospaces",
			wantErr: true,
		},
		{
			name:    "only one space",
			line:    "only one",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileAndPos, stmts, count, err := splitCoverLineParts(tt.line)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fileAndPos != tt.wantFile {
				t.Errorf("fileAndPos = %q, want %q", fileAndPos, tt.wantFile)
			}
			if stmts != tt.wantStmts {
				t.Errorf("stmts = %q, want %q", stmts, tt.wantStmts)
			}
			if count != tt.wantCount {
				t.Errorf("count = %q, want %q", count, tt.wantCount)
			}
		})
	}
}
