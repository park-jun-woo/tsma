package index

import "testing"

func TestLastNonEmptyBeforeLine(t *testing.T) {
	tests := []struct {
		name         string
		currentLine  int
		lastNonEmpty int
		want         int
	}{
		{"lastNonEmpty before current", 10, 6, 6},
		{"lastNonEmpty equals current falls back", 5, 5, 4},
		{"lastNonEmpty after current falls back", 5, 9, 4},
		{"current line is one", 1, 1, 1},
		{"current line is one, lastNonEmpty after", 1, 3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lastNonEmptyBeforeLine(tt.currentLine, tt.lastNonEmpty)
			if got != tt.want {
				t.Errorf("lastNonEmptyBeforeLine(%d, %d) = %d, want %d",
					tt.currentLine, tt.lastNonEmpty, got, tt.want)
			}
		})
	}
}
