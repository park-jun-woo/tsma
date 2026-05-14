package coverage

import "testing"

func TestParseCoverPositions(t *testing.T) {
	tests := []struct {
		name      string
		positions string
		wantStart int
		wantEnd   int
		wantErr   bool
	}{
		{
			name:      "valid positions",
			positions: "10.2,20.5",
			wantStart: 10,
			wantEnd:   20,
		},
		{
			name:      "single digit positions",
			positions: "1.1,2.2",
			wantStart: 1,
			wantEnd:   2,
		},
		{
			name:    "no comma",
			positions: "10.2",
			wantErr: true,
		},
		{
			name:    "too many commas",
			positions: "10.2,20.5,30.1",
			wantErr: true,
		},
		{
			name:    "no dot in start",
			positions: "10,20.5",
			wantErr: true,
		},
		{
			name:    "no dot in end",
			positions: "10.2,20",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &coverBlock{}
			err := parseCoverPositions(tt.positions, b)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if b.startLine != tt.wantStart {
				t.Errorf("startLine = %d, want %d", b.startLine, tt.wantStart)
			}
			if b.endLine != tt.wantEnd {
				t.Errorf("endLine = %d, want %d", b.endLine, tt.wantEnd)
			}
		})
	}
}
