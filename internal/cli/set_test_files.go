//ff:func feature=cli type=helper control=sequence
//ff:what Records a TestMatch onto a function as TestFiles plus a representative TestFile
package cli

import (
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// setTestFiles records a match onto a function: the full file set in TestFiles
// plus the representative TestFile (the first file) for display and back-compat.
// A match with no files is a no-op.
func setTestFiles(fn *model.Function, tm match.TestMatch) {
	if len(tm.Files) == 0 {
		return
	}
	fn.TestFiles = tm.Files
	fn.TestFile = tm.Files[0]
}
