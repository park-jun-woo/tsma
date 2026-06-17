//ff:func feature=smell type=test
//ff:what linknameInGroup 단위테스트: 지시자 그룹(//go:linkname → 발화)과 일반 산문 그룹(무발화)을 직접 호출해 HasPrefix 양/음 분기를 모두 덮는다.

package smell

import "testing"

func TestLinknameInGroup_DirectiveAndProse(t *testing.T) {
	// Two comment groups: a //go:linkname directive (fires) and a prose comment
	// whose Text starts with "// " (never fires).
	src := `package p

//go:linkname privateFn runtime.someInternal
func privateFn()

// prose comment mentioning go:linkname inline
func other() {}
`
	file, fset, path := parseDetectorSrc(t, src)

	var matched int
	for _, g := range file.Comments {
		matched += len(linknameInGroup(g, fset, path))
	}
	if matched != 1 {
		t.Fatalf("matched = %d, want 1 (only the directive group fires)", matched)
	}
}
