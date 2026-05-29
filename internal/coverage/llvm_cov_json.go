//ff:type feature=coverage type=model
//ff:what Represents the top-level structure of cargo llvm-cov --json (LLVM coverage export) output
package coverage

// llvmCovJSON represents the top-level LLVM coverage export document.
type llvmCovJSON struct {
	Data []llvmCovData `json:"data"`
}
