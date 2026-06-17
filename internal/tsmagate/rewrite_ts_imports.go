//ff:func feature=gate type=helper control=sequence lang=typescript
//ff:what rewriteTSImports: re-points the relative module specifiers in a generated test so that, when the test is placed in the isolation dir (.tsma/test) instead of beside the source, `import { foo } from "./mod"` still resolves to the real module. Each "./" or "../" specifier is resolved against the source dir, then re-expressed relative to the backing dir. Bare/package specifiers (no leading dot) are left untouched. This is what makes the .tsma/test isolation import the real module by path (plan §6).
package tsmagate

import (
	"path/filepath"
	"regexp"
	"strings"
)

// tsImportSpecifier matches the module string of an `import ... from "x"`,
// `export ... from "x"`, or dynamic `import("x")` — capturing the quote and the
// specifier so only the path is rewritten.
var tsImportSpecifier = regexp.MustCompile(`(from\s*|import\s*\(\s*)(['"])(\.[^'"]*)(['"])`)

// rewriteTSImports rewrites every relative import specifier in src so it resolves
// from backingDirAbs to the module it named relative to sourceDirAbs. Absolute
// and bare package specifiers are left unchanged.
func rewriteTSImports(src, sourceDirAbs, backingDirAbs string) string {
	return tsImportSpecifier.ReplaceAllStringFunc(src, func(m string) string {
		groups := tsImportSpecifier.FindStringSubmatch(m)
		prefix, openQ, spec, closeQ := groups[1], groups[2], groups[3], groups[4]
		target := filepath.Join(sourceDirAbs, spec)
		rel, err := filepath.Rel(backingDirAbs, target)
		if err != nil {
			return m
		}
		rel = filepath.ToSlash(rel)
		if !strings.HasPrefix(rel, ".") {
			rel = "./" + rel
		}
		return prefix + openQ + rel + closeQ
	})
}
