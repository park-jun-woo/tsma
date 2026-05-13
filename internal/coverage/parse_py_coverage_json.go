//ff:func feature=coverage type=implementation control=sequence
//ff:what Reads and parses the coverage.json file from Python coverage tools
package coverage

import (
	"encoding/json"
	"fmt"
	"os"
)

// parsePyCoverageJSON reads and parses the coverage.json file.
func parsePyCoverageJSON(path string) (*pyCoverageJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var covData pyCoverageJSON
	if err := json.Unmarshal(data, &covData); err != nil {
		return nil, fmt.Errorf("unmarshal coverage json: %w", err)
	}

	return &covData, nil
}
