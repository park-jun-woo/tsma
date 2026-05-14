package coverage

import "testing"

func TestOverlapsTable(t *testing.T) {
	tests := []struct {
		name       string
		blockFile  string
		targetFile string
		blockStart int
		blockEnd   int
		funcStart  int
		funcEnd    int
		want       bool
	}{
		{
			name:       "block inside func range",
			blockFile:  "github.com/example/pkg/handler.go",
			targetFile: "pkg/handler.go",
			blockStart: 15, blockEnd: 20,
			funcStart: 10, funcEnd: 25,
			want: true,
		},
		{
			name:       "block before func range",
			blockFile:  "github.com/example/pkg/handler.go",
			targetFile: "pkg/handler.go",
			blockStart: 5, blockEnd: 8,
			funcStart: 10, funcEnd: 25,
			want: false,
		},
		{
			name:       "different file",
			blockFile:  "github.com/example/pkg/service.go",
			targetFile: "pkg/handler.go",
			blockStart: 15, blockEnd: 20,
			funcStart: 10, funcEnd: 25,
			want: false,
		},
		{
			name:       "block starts at func start",
			blockFile:  "github.com/example/pkg/handler.go",
			targetFile: "pkg/handler.go",
			blockStart: 10, blockEnd: 12,
			funcStart: 10, funcEnd: 25,
			want: true,
		},
		{
			name:       "block starts at func end",
			blockFile:  "github.com/example/pkg/handler.go",
			targetFile: "pkg/handler.go",
			blockStart: 25, blockEnd: 30,
			funcStart: 10, funcEnd: 25,
			want: true,
		},
		{
			name:       "block starts after func end",
			blockFile:  "github.com/example/pkg/handler.go",
			targetFile: "pkg/handler.go",
			blockStart: 26, blockEnd: 30,
			funcStart: 10, funcEnd: 25,
			want: false,
		},
		{
			name:       "exact file match",
			blockFile:  "handler.go",
			targetFile: "handler.go",
			blockStart: 10, blockEnd: 20,
			funcStart: 5, funcEnd: 30,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := overlaps(tt.blockFile, tt.targetFile, tt.blockStart, tt.blockEnd, tt.funcStart, tt.funcEnd)
			if got != tt.want {
				t.Errorf("overlaps(%q, %q, %d, %d, %d, %d) = %v, want %v",
					tt.blockFile, tt.targetFile, tt.blockStart, tt.blockEnd, tt.funcStart, tt.funcEnd, got, tt.want)
			}
		})
	}
}
