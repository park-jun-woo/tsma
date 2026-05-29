package coverage

import (
	"encoding/json"
	"testing"
)

func TestLLVMSegmentUnmarshal(t *testing.T) {
	var s llvmSegment
	if err := json.Unmarshal([]byte(`[3, 9, 2, true, true, false]`), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.Line != 3 || s.Col != 9 || s.Count != 2 || !s.HasCount || !s.IsRegionEntry || s.IsGapRegion {
		t.Errorf("decoded = %+v", s)
	}
}

func TestLLVMSegmentUnmarshalShort(t *testing.T) {
	var s llvmSegment
	if err := json.Unmarshal([]byte(`[1, 2]`), &s); err == nil {
		t.Error("expected error for too-short segment array")
	}
}

func TestLLVMSegmentUnmarshalFiveElements(t *testing.T) {
	// Exactly 5 elements: the optional 6th-element branch is skipped and
	// IsGapRegion stays false.
	var s llvmSegment
	if err := json.Unmarshal([]byte(`[7, 4, 1, true, false]`), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.Line != 7 || s.Col != 4 || s.Count != 1 || !s.HasCount || s.IsRegionEntry || s.IsGapRegion {
		t.Errorf("decoded = %+v", s)
	}
}

func TestLLVMSegmentUnmarshalNotArray(t *testing.T) {
	var s llvmSegment
	// Outer unmarshal into []json.RawMessage fails.
	if err := json.Unmarshal([]byte(`{"x":1}`), &s); err == nil {
		t.Error("expected error for non-array segment JSON")
	}
}

func TestLLVMSegmentUnmarshalBadFields(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"badLine", `["x", 9, 2, true, true, false]`},
		{"badCol", `[3, "x", 2, true, true, false]`},
		{"badCount", `[3, 9, "x", true, true, false]`},
		{"badHasCount", `[3, 9, 2, "x", true, false]`},
		{"badIsRegionEntry", `[3, 9, 2, true, "x", false]`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s llvmSegment
			if err := json.Unmarshal([]byte(tc.in), &s); err == nil {
				t.Errorf("expected error for %s", tc.name)
			}
		})
	}
}
