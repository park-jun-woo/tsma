//ff:func feature=coverage type=implementation control=sequence lang=csharp
//ff:what Reads and parses a Cobertura cobertura.xml report into flattened per-file coverage
package coverage

import (
	"encoding/xml"
	"fmt"
	"os"
)

// parseCobertura reads the Cobertura XML report at path (as produced by coverlet
// via `dotnet test --collect:"XPlat Code Coverage"`) and returns its coverage
// data flattened to one entry per source file keyed by class filename.
func parseCobertura(path string) (*csCoverage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var report coberturaReport
	if err := xml.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("decode cobertura xml: %w", err)
	}

	return flattenCobertura(&report), nil
}
