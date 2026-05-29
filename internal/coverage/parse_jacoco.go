//ff:func feature=coverage type=implementation control=sequence
//ff:what Reads and parses a JaCoCo jacoco.xml report into flattened per-file coverage
package coverage

import (
	"encoding/xml"
	"fmt"
	"os"
)

// parseJacoco reads the JaCoCo XML report at path and returns its coverage data
// flattened to one entry per source file (with package-qualified paths). The
// report's external DTD reference is ignored; Go's decoder does not fetch it.
func parseJacoco(path string) (*jacocoCoverage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var report jacocoReport
	if err := xml.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("decode jacoco xml: %w", err)
	}

	return flattenJacoco(&report), nil
}
