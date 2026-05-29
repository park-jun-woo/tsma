package index

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestClosePrevTSEndLineAtEOF(t *testing.T) {
	t.Run("empty slice is a no-op", func(t *testing.T) {
		var fns []model.Function
		closePrevTSEndLineAtEOF(fns, 42)
		if len(fns) != 0 {
			t.Errorf("expected empty slice unchanged")
		}
	})

	t.Run("closes when EndLine equals StartLine", func(t *testing.T) {
		fns := []model.Function{
			{Name: "A", StartLine: 5, EndLine: 5},
		}
		closePrevTSEndLineAtEOF(fns, 30)
		if fns[0].EndLine != 30 {
			t.Errorf("expected EndLine=30, got %d", fns[0].EndLine)
		}
	})

	t.Run("leaves already-closed function", func(t *testing.T) {
		fns := []model.Function{
			{Name: "A", StartLine: 5, EndLine: 20},
		}
		closePrevTSEndLineAtEOF(fns, 30)
		if fns[0].EndLine != 20 {
			t.Errorf("expected EndLine to remain 20, got %d", fns[0].EndLine)
		}
	})
}
