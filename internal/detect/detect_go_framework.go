//ff:func feature=detect type=implementation control=sequence
//ff:what Reads go.mod and detects the Go web framework by import path
package detect

import (
	"os"
	"path/filepath"
)

func detectGoFramework(projectRoot string) string {
	data, err := os.ReadFile(filepath.Join(projectRoot, "go.mod"))
	if err != nil {
		return "unknown"
	}
	content := string(data)
	if containsImport(content, "github.com/gin-gonic/gin") {
		return "gin"
	}
	if containsImport(content, "github.com/labstack/echo") {
		return "echo"
	}
	if containsImport(content, "github.com/go-chi/chi") {
		return "chi"
	}
	return "unknown"
}
