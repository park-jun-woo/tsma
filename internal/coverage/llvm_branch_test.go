package coverage

import (
	"encoding/json"
	"testing"
)

func TestLLVMBranchUnmarshal(t *testing.T) {
	var b llvmBranch
	if err := json.Unmarshal([]byte(`[3, 9, 3, 19, 2, 3, 0, 0, 0]`), &b); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if b.LineStart != 3 || b.ExecCount != 2 || b.FalseExecCount != 3 {
		t.Errorf("decoded = %+v", b)
	}
}

func TestLLVMBranchUnmarshalShort(t *testing.T) {
	var b llvmBranch
	if err := json.Unmarshal([]byte(`[1, 2, 3]`), &b); err == nil {
		t.Error("expected error for too-short branch array")
	}
}

func TestLLVMBranchUnmarshalInvalidJSON(t *testing.T) {
	var b llvmBranch
	// Not an array of ints -> the inner json.Unmarshal into []int fails,
	// exercising the early error-return branch.
	if err := json.Unmarshal([]byte(`{"not":"an array"}`), &b); err == nil {
		t.Error("expected error for non-array branch JSON")
	}
}

func TestLLVMBranchUnmarshalWrongElemType(t *testing.T) {
	var b llvmBranch
	// Array with a string element cannot decode into []int.
	if err := json.Unmarshal([]byte(`[1, "x", 3, 4, 5, 6]`), &b); err == nil {
		t.Error("expected error for non-int branch element")
	}
}
