//ff:func feature=cli type=helper control=sequence
//ff:what Finds a function by name, reporting ambiguity if multiple matches exist
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/tsma/internal/model"
)

// findFunctionOrAmbiguous finds a function by name. If not found, checks for
// ambiguous matches and prints them. Returns the function or an error.
func findFunctionOrAmbiguous(sess *model.Session, funcName string) (*model.Function, error) {
	fn := sess.FindFunction(funcName)
	if fn != nil {
		return fn, nil
	}

	matches := sess.FindAmbiguous(funcName)
	if len(matches) > 1 {
		fmt.Fprintf(os.Stderr, "ambiguous function name %q — %d matches:\n", funcName, len(matches))
		for _, m := range matches {
			fmt.Fprintf(os.Stderr, "  %s  %s:%d-%d\n", m.QualifiedName, m.File, m.StartLine, m.EndLine)
		}
		return nil, fmt.Errorf("use the qualified name to disambiguate")
	}
	return nil, fmt.Errorf("function not found: %s", funcName)
}
