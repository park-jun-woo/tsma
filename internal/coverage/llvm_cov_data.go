//ff:type feature=coverage type=model
//ff:what Represents a single data block of the LLVM coverage export containing per-file coverage
package coverage

// llvmCovData is a single export entry containing per-file coverage.
type llvmCovData struct {
	Files []llvmCovFile `json:"files"`
}
