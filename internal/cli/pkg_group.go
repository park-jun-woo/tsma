//ff:type feature=cli type=model
//ff:what Holds functions grouped by their package directory for batch test execution
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// pkgGroup holds functions grouped by their package directory.
type pkgGroup struct {
	pkgDir   string
	testFile string
	funcs    []*model.Function
}
