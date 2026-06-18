//ff:func feature=detect type=helper control=iteration dimension=1 lang=python
//ff:what Reports whether a project-local venv carries a pytest executable
package detect

import (
	"path/filepath"
)

// probePytestVenv reports whether a project-local virtualenv ships a pytest
// executable at the standard POSIX location (<venv>/bin/pytest) for the .venv or
// venv directory. This is the last-resort signal, used only when static
// declaration signals miss: findPython() returns the system python3 (not the
// venv), so a subprocess `python -m pytest --version` would not see a venv-only
// pytest — checking the file directly is both cheaper and more accurate.
// Windows (Scripts/) layout is out of scope (Phase006 §4).
func probePytestVenv(projectRoot string) bool {
	for _, venv := range []string{".venv", "venv"} {
		if fileExists(filepath.Join(projectRoot, venv, "bin", "pytest")) {
			return true
		}
	}
	return false
}
