//ff:type feature=cli type=model
//ff:what Groups one Go package's matched functions and the union of their test funcs for a single test run
package cli

import (
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// goPkgGroup collects the functions of one package that have a matched test,
// along with the union of their test functions, so the package can be measured
// with a single `go test` run.
type goPkgGroup struct {
	funcs     []*model.Function
	matches   map[*model.Function]match.TestMatch
	testFuncs []string
}
