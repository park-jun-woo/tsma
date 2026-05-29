//ff:type feature=coverage type=model
//ff:what Represents per-file LLVM coverage holding line segments and branch regions
package coverage

// llvmCovFile holds line segments and branch regions for one source file.
type llvmCovFile struct {
	Filename string        `json:"filename"`
	Segments []llvmSegment `json:"segments"`
	Branches []llvmBranch  `json:"branches"`
}
