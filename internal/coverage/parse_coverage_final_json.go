//ff:func feature=coverage type=implementation control=sequence
//ff:what Reads coverage-final.json from the coverage directory and unmarshals istanbul format
package coverage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// parseCoverageFinalJSON reads and parses the coverage-final.json file from the coverage directory.
func parseCoverageFinalJSON(coverDir string) (map[string]coverageFinalEntry, error) {
	coverFile := filepath.Join(coverDir, "coverage-final.json")
	data, err := os.ReadFile(coverFile)
	if err != nil {
		return nil, fmt.Errorf("read coverage file: %w", err)
	}

	var result map[string]coverageFinalEntry
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal coverage JSON: %w", err)
	}

	return result, nil
}
