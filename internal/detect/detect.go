//ff:func feature=detect type=factory control=sequence
//ff:what Identifies the project language and framework from marker files
package detect

import (
	"os"
	"path/filepath"
)

// Detect identifies the project language and framework.
// MVP: Go only.
func Detect(projectRoot string) (*LangFramework, error) {
	// Go
	if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
		fw := detectGoFramework(projectRoot)
		return &LangFramework{Lang: "go", Framework: fw}, nil
	}

	// TypeScript / JavaScript
	if _, err := os.Stat(filepath.Join(projectRoot, "package.json")); err == nil {
		fw := detectTSFramework(projectRoot)
		return &LangFramework{Lang: "typescript", Framework: fw}, nil
	}

	// Python
	if _, err := os.Stat(filepath.Join(projectRoot, "pyproject.toml")); err == nil {
		fw := detectPyFramework(projectRoot)
		return &LangFramework{Lang: "python", Framework: fw}, nil
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "requirements.txt")); err == nil {
		fw := detectPyFramework(projectRoot)
		return &LangFramework{Lang: "python", Framework: fw}, nil
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "setup.py")); err == nil {
		fw := detectPyFramework(projectRoot)
		return &LangFramework{Lang: "python", Framework: fw}, nil
	}

	return nil, ErrUnsupportedLanguage
}
