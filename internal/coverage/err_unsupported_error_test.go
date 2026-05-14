package coverage

import (
	"strings"
	"testing"
)

func TestErrUnsupportedError(t *testing.T) {
	tests := []struct {
		lang     string
		wantSub  string
	}{
		{"rust", "rust"},
		{"java", "java"},
		{"ruby", "ruby"},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			err := &ErrUnsupported{Lang: tt.lang}
			msg := err.Error()
			if !strings.Contains(msg, tt.wantSub) {
				t.Errorf("Error() = %q, want to contain %q", msg, tt.wantSub)
			}
			if !strings.Contains(msg, "not implemented") {
				t.Errorf("Error() = %q, want to contain 'not implemented'", msg)
			}
		})
	}
}

func TestErrUnsupportedImplementsError(t *testing.T) {
	var err error = &ErrUnsupported{Lang: "test"}
	if err.Error() == "" {
		t.Error("Error() should return non-empty string")
	}
}
