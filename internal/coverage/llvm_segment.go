//ff:type feature=coverage type=model
//ff:what Represents one LLVM coverage line segment decoded from its positional JSON array
package coverage

import (
	"encoding/json"
	"fmt"
)

// llvmSegment is one entry of a file's "segments" array. LLVM encodes each
// segment as the tuple [line, col, count, hasCount, isRegionEntry, isGapRegion].
type llvmSegment struct {
	Line          int
	Col           int
	Count         int
	HasCount      bool
	IsRegionEntry bool
	IsGapRegion   bool
}

// UnmarshalJSON decodes the positional segment array into named fields.
func (s *llvmSegment) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) < 5 {
		return fmt.Errorf("llvm segment: expected at least 5 elements, got %d", len(raw))
	}
	if err := json.Unmarshal(raw[0], &s.Line); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[1], &s.Col); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[2], &s.Count); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[3], &s.HasCount); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[4], &s.IsRegionEntry); err != nil {
		return err
	}
	if len(raw) >= 6 {
		_ = json.Unmarshal(raw[5], &s.IsGapRegion)
	}
	return nil
}
