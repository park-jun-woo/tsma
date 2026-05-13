//ff:func feature=detect type=implementation control=iteration dimension=1
//ff:what Reads Python dependency files and detects the web framework
package detect

import (
	"os"
	"path/filepath"
)

func detectPyFramework(projectRoot string) string {
	for _, name := range []string{"pyproject.toml", "requirements.txt", "setup.py"} {
		data, err := os.ReadFile(filepath.Join(projectRoot, name))
		if err != nil {
			continue
		}
		content := string(data)
		if containsImport(content, "fastapi") {
			return "fastapi"
		}
		if containsImport(content, "django") {
			return "django"
		}
	}
	return "unknown"
}
