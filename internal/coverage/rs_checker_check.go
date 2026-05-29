//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs cargo llvm-cov and computes per-function Rust coverage from its JSON output
package coverage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Check runs cargo llvm-cov for the project and maps the JSON export onto the
// given function's line range. Requires a cargo + cargo-llvm-cov toolchain.
func (c *RsChecker) Check(projectRoot, testFile string, fn *model.Function) (*Report, error) {
	cargo, err := findCargo()
	if err != nil {
		return nil, err
	}

	tmDir := filepath.Join(projectRoot, ".tsma")
	if err := os.MkdirAll(tmDir, 0o755); err != nil {
		return nil, fmt.Errorf("create .tsma dir: %w", err)
	}

	coverJSONPath := filepath.Join(tmDir, "llvm-cov.json")

	if err := runCargoLLVMCov(cargo, projectRoot, coverJSONPath); err != nil {
		return nil, fmt.Errorf("rust coverage failed: %w", err)
	}

	cov, err := parseLLVMCov(coverJSONPath)
	if err != nil {
		return nil, fmt.Errorf("parse llvm-cov json: %w", err)
	}

	ranges := collectRsRanges(fn)

	return buildRsReport(ranges, cov, projectRoot), nil
}
