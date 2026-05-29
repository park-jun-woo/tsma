package index

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestClosePrevTSEndLine(t *testing.T) {
	t.Run("empty slice is a no-op", func(t *testing.T) {
		var fns []model.Function
		closePrevTSEndLine(fns, 10, 8)
		if len(fns) != 0 {
			t.Errorf("expected empty slice unchanged")
		}
	})

	t.Run("closes when EndLine equals StartLine", func(t *testing.T) {
		fns := []model.Function{
			{Name: "A", StartLine: 3, EndLine: 3},
		}
		// lastNonEmpty (8) is not < lineNum (5) -> fallback lineNum-1 = 4.
		closePrevTSEndLine(fns, 5, 8)
		if fns[0].EndLine != 4 {
			t.Errorf("expected EndLine=4, got %d", fns[0].EndLine)
		}
	})

	t.Run("uses lastNonEmpty when before current line", func(t *testing.T) {
		fns := []model.Function{
			{Name: "A", StartLine: 3, EndLine: 3},
		}
		// lastNonEmpty (6) < lineNum (10) -> use lastNonEmpty.
		closePrevTSEndLine(fns, 10, 6)
		if fns[0].EndLine != 6 {
			t.Errorf("expected EndLine=6, got %d", fns[0].EndLine)
		}
	})

	t.Run("leaves already-closed function", func(t *testing.T) {
		fns := []model.Function{
			{Name: "A", StartLine: 3, EndLine: 12},
		}
		closePrevTSEndLine(fns, 20, 18)
		if fns[0].EndLine != 12 {
			t.Errorf("expected EndLine to remain 12, got %d", fns[0].EndLine)
		}
	})
}
