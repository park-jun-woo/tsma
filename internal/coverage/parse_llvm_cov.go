//ff:func feature=coverage type=implementation control=sequence
//ff:what Reads and parses a cargo llvm-cov --json output file
package coverage

import (
	"encoding/json"
	"fmt"
	"os"
)

// parseLLVMCov reads and parses the llvm-cov JSON export at path.
func parseLLVMCov(path string) (*llvmCovJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cov llvmCovJSON
	if err := json.Unmarshal(data, &cov); err != nil {
		return nil, fmt.Errorf("unmarshal llvm-cov json: %w", err)
	}

	return &cov, nil
}
