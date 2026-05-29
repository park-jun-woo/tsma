package index

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestClosePrevEndLine(t *testing.T) {
	t.Run("empty slice is a no-op", func(t *testing.T) {
		var fns []model.Function
		closePrevEndLine(fns, 42) // must not panic
		if len(fns) != 0 {
			t.Errorf("expected slice unchanged, got len %d", len(fns))
		}
	})

	t.Run("sets EndLine when still zero", func(t *testing.T) {
		fns := []model.Function{
			{Name: "A", StartLine: 1, EndLine: 5},
			{Name: "B", StartLine: 6, EndLine: 0},
		}
		closePrevEndLine(fns, 10)
		if fns[1].EndLine != 10 {
			t.Errorf("expected EndLine=10, got %d", fns[1].EndLine)
		}
		// Earlier entries must be untouched.
		if fns[0].EndLine != 5 {
			t.Errorf("expected first EndLine untouched at 5, got %d", fns[0].EndLine)
		}
	})

	t.Run("leaves EndLine when already set", func(t *testing.T) {
		fns := []model.Function{
			{Name: "B", StartLine: 6, EndLine: 8},
		}
		closePrevEndLine(fns, 99)
		if fns[0].EndLine != 8 {
			t.Errorf("expected EndLine to stay 8, got %d", fns[0].EndLine)
		}
	})
}
