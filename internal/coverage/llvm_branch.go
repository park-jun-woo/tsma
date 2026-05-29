//ff:type feature=coverage type=model
//ff:what Represents one LLVM coverage branch region decoded from its positional JSON array
package coverage

import (
	"encoding/json"
	"fmt"
)

// llvmBranch is one entry of a file's "branches" array. LLVM encodes each
// branch region as the tuple
// [lineStart, colStart, lineEnd, colEnd, execCount, falseExecCount, fileID, expandedFileID, regionKind].
type llvmBranch struct {
	LineStart      int
	ExecCount      int
	FalseExecCount int
}

// UnmarshalJSON decodes the positional branch array into the fields needed for
// per-function branch coverage mapping.
func (b *llvmBranch) UnmarshalJSON(data []byte) error {
	var raw []int
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) < 6 {
		return fmt.Errorf("llvm branch: expected at least 6 elements, got %d", len(raw))
	}
	b.LineStart = raw[0]
	b.ExecCount = raw[4]
	b.FalseExecCount = raw[5]
	return nil
}
