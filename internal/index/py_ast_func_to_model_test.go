package index

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestPyAstFuncToModel maps single ast entries to model.Function, covering the
// exported heuristic (uppercase + no leading underscore), receiver qualified
// names, and verbatim line-range passthrough.
func TestPyAstFuncToModel(t *testing.T) {
	tests := []struct {
		name      string
		f         pyAstFunc
		relDir    string
		relPath   string
		wantQual  string
		wantExp   bool
		wantRecv  string
		wantStart int
		wantEnd   int
	}{
		{
			name:     "lowercase top-level is unexported",
			f:        pyAstFunc{Name: "classify", StartLine: 1, EndLine: 4},
			relDir:   "src",
			relPath:  "src/calc.py",
			wantQual: "src.classify",
			wantExp:  false, wantRecv: "", wantStart: 1, wantEnd: 4,
		},
		{
			name:     "uppercase method is exported with receiver",
			f:        pyAstFunc{Name: "Add", Receiver: "Calculator", StartLine: 12, EndLine: 15},
			relDir:   "src",
			relPath:  "src/calc.py",
			wantQual: "src.Calculator.Add",
			wantExp:  true, wantRecv: "Calculator", wantStart: 12, wantEnd: 15,
		},
		{
			name:     "leading underscore is unexported",
			f:        pyAstFunc{Name: "_inner", StartLine: 13, EndLine: 14},
			relDir:   "src",
			relPath:  "src/calc.py",
			wantQual: "src._inner",
			wantExp:  false, wantRecv: "", wantStart: 13, wantEnd: 14,
		},
		{
			name:     "empty name is unexported",
			f:        pyAstFunc{Name: "", StartLine: 1, EndLine: 1},
			relDir:   "src",
			relPath:  "src/calc.py",
			wantExp:  false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := pyAstFuncToModel(tc.f, tc.relDir, tc.relPath)
			if tc.wantQual != "" && got.QualifiedName != tc.wantQual {
				t.Errorf("QualifiedName = %q, want %q", got.QualifiedName, tc.wantQual)
			}
			if got.Exported != tc.wantExp {
				t.Errorf("Exported = %v, want %v", got.Exported, tc.wantExp)
			}
			if got.Receiver != tc.wantRecv {
				t.Errorf("Receiver = %q, want %q", got.Receiver, tc.wantRecv)
			}
			if got.Name != tc.f.Name {
				t.Errorf("Name = %q, want %q", got.Name, tc.f.Name)
			}
			if got.File != tc.relPath {
				t.Errorf("File = %q, want %q", got.File, tc.relPath)
			}
			if tc.wantStart != 0 && (got.StartLine != tc.wantStart || got.EndLine != tc.wantEnd) {
				t.Errorf("range = %d-%d, want %d-%d", got.StartLine, got.EndLine, tc.wantStart, tc.wantEnd)
			}
			if got.Status != model.StatusTodo {
				t.Errorf("Status = %v, want StatusTodo", got.Status)
			}
		})
	}
}
