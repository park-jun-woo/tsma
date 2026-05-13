//ff:func feature=detect type=factory control=sequence
//ff:what Identifies the project language from marker files
package detect

import (
	"os"
	"path/filepath"
)

// Detect identifies the project language.
func Detect(projectRoot string) (*LangFramework, error) {
	// Go
	if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
		return &LangFramework{Lang: "go"}, nil
	}

	// TypeScript / JavaScript
	if _, err := os.Stat(filepath.Join(projectRoot, "package.json")); err == nil {
		return &LangFramework{Lang: "typescript"}, nil
	}

	// Python
	if _, err := os.Stat(filepath.Join(projectRoot, "pyproject.toml")); err == nil {
		return &LangFramework{Lang: "python"}, nil
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "requirements.txt")); err == nil {
		return &LangFramework{Lang: "python"}, nil
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "setup.py")); err == nil {
		return &LangFramework{Lang: "python"}, nil
	}

	return nil, ErrUnsupportedLanguage
}
