//ff:func feature=match type=helper control=iteration dimension=1
//ff:what Finds a misnamed test_<base>_test.go variant when the canonical Go test file is absent
package match

import (
	"os"
	"path/filepath"
)

// FindMisnamedTest checks whether a non-canonical test_ prefixed Go test file
// (test_<base>_test.go) exists in the source directory while the canonical
// <base>_test.go is missing. Only Go sources are handled; any other lang
// (including python) always returns found=false.
func FindMisnamedTest(projectRoot, lang, sourceFile string) (misnamedRel, canonicalRel string, found bool) {
	// The canonical <base>_test.go path comes from the shared formula; "" means
	// a non-Go lang or a non-.go source — neither has a misnamed variant.
	canonicalRel = CanonicalTestPath(lang, sourceFile)
	if canonicalRel == "" {
		return "", "", false
	}

	srcDir := filepath.Dir(sourceFile)
	canonicalBase := filepath.Base(canonicalRel)

	// Canonical test file already exists -> normal, do not intervene.
	if _, err := os.Stat(filepath.Join(projectRoot, canonicalRel)); err == nil {
		return "", "", false
	}

	// Only candidate: test_<base>_test.go (Python-convention prefix mixed in).
	// test_<base>.go is intentionally excluded (it is a production source per isGoSource).
	for _, candidate := range []string{"test_" + canonicalBase} {
		candidateRel := filepath.Join(srcDir, candidate)
		if _, err := os.Stat(filepath.Join(projectRoot, candidateRel)); err == nil {
			return candidateRel, canonicalRel, true
		}
	}

	return "", "", false
}
