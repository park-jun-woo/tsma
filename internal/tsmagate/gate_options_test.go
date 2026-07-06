//ff:func feature=gate type=test
//ff:what GateOptions 단위테스트: go:embed된 gate.md 원문(Src 비어있지 않음·헤더 포함), 로드 에러 표시용 Path, 명시 Case "제출 통과", Registry 비nil, ground 미사용 계약(Provider/Resolver nil)을 고정한다.

package tsmagate

import (
	"strings"
	"testing"
)

// TestGateOptions pins the GateOptions wiring: Src carries the embedded
// gate.md (non-empty, with the doc header), Path names the source file for
// load errors, Case is explicitly "제출 통과", Registry is bound, and the
// no-ground contract holds (Provider and Resolver are nil — all tier-0).
func TestGateOptions(t *testing.T) {
	opts := GateOptions()
	if opts == nil {
		t.Fatal("GateOptions() = nil")
	}
	if len(opts.Src) == 0 {
		t.Error("Src is empty — gate.md was not embedded")
	}
	if !strings.Contains(string(opts.Src), "# tsma 게이트") {
		t.Error("Src does not contain the gate.md header")
	}
	if opts.Path != "internal/tsmagate/gate.md" {
		t.Errorf("Path = %q, want internal/tsmagate/gate.md", opts.Path)
	}
	if opts.Case != "제출 통과" {
		t.Errorf("Case = %q, want 제출 통과", opts.Case)
	}
	if opts.Registry == nil {
		t.Error("Registry is nil")
	}
	if opts.Provider != nil {
		t.Error("Provider must stay nil (no ground use)")
	}
	if opts.Resolver != nil {
		t.Error("Resolver must stay nil (no network seam)")
	}
}
