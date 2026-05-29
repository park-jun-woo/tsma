//ff:func feature=coverage type=helper control=sequence lang=csharp
//ff:what Locates the Cobertura XML report produced by coverlet under a results directory
package coverage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// findCoberturaReport searches resultsDir recursively for the Cobertura report
// emitted by coverlet's XPlat collector, which writes
// <resultsDir>/<guid>/coverage.cobertura.xml. The first matching *.cobertura.xml
// is returned; an error is returned when none is found.
func findCoberturaReport(resultsDir string) (string, error) {
	var found string
	err := filepath.Walk(resultsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if found != "" {
			return filepath.SkipDir
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "cobertura.xml") {
			found = path
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if found == "" {
		return "", fmt.Errorf("no cobertura.xml report found under %s", resultsDir)
	}
	return found, nil
}
