package coverage

import "testing"

func TestParseCoverLineValid(t *testing.T) {
	line := "github.com/example/pkg/handler.go:10.2,20.5 3 1"
	b, err := parseCoverLine(line)
	if err != nil {
		t.Fatalf("parseCoverLine: %v", err)
	}
	if b.file != "github.com/example/pkg/handler.go" {
		t.Errorf("file = %q", b.file)
	}
	if b.startLine != 10 || b.startCol != 2 {
		t.Errorf("start = %d.%d, want 10.2", b.startLine, b.startCol)
	}
	if b.endLine != 20 || b.endCol != 5 {
		t.Errorf("end = %d.%d, want 20.5", b.endLine, b.endCol)
	}
	if b.stmts != 3 {
		t.Errorf("stmts = %d, want 3", b.stmts)
	}
	if b.count != 1 {
		t.Errorf("count = %d, want 1", b.count)
	}
}

func TestParseCoverLineZeroCount(t *testing.T) {
	line := "pkg/main.go:1.1,5.5 2 0"
	b, err := parseCoverLine(line)
	if err != nil {
		t.Fatalf("parseCoverLine: %v", err)
	}
	if b.count != 0 {
		t.Errorf("count = %d, want 0", b.count)
	}
}

func TestParseCoverLineInvalid(t *testing.T) {
	tests := []struct {
		name string
		line string
	}{
		{"no spaces", "nospaces"},
		{"only one space", "one space"},
		// Passes splitCoverLineParts but the file/pos part has no colon ->
		// extractCoverFile returns an error.
		{"no colon in file part", "nocolon 3 1"},
		// Valid file:pos prefix, but positions contain an extra comma so
		// parseCoverPositions fails (3 parts instead of 2).
		{"bad positions", "file.go:1.1,2.2,3.3 3 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseCoverLine(tt.line)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
