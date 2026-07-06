//ff:func feature=gate type=test
//ff:what mergeHeaders 단위테스트: 기존-먼저 순서 안정성, 동일 라인 dedup, 빈 줄
// 연속 붕괴(단일 구분자 유지), out이 비어 있을 때의 빈 줄 유지까지 전 분기를 덮는다.
package tsmagate

import (
	"reflect"
	"testing"
)

func TestMergeHeaders(t *testing.T) {
	existing := []string{"", `import "a"`, "", "", `import "b"`}
	incoming := []string{`import "a"`, `import "c"`}
	got := mergeHeaders(existing, incoming)
	want := []string{"", `import "a"`, "", `import "b"`, `import "c"`}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mergeHeaders = %q, want %q", got, want)
	}
}
