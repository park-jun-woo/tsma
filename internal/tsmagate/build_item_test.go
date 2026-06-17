//ff:func feature=gate type=test
//ff:what buildItem 단위테스트: 함수 하나를 quest.Item으로 변환하는 성공 경로를 덮는다. 런타임 필드(Status/Attempt/CoveragePct/TestMtime/FailOutput) 제거와 Key=QualifiedName·State=TODO·payload(lang/root/fn) 스냅샷을 검증한다.

package tsmagate

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/model"
)

func TestBuildItem_StripsAndSnapshots(t *testing.T) {
	fn := model.Function{
		QualifiedName: "pkg.Foo",
		Name:          "Foo",
		Status:        "DONE",
		Attempt:       3,
		CoveragePct:   50,
		TestMtime:     "t",
		FailOutput:    "boom",
	}
	it, err := buildItem("go", "/root", fn)
	if err != nil {
		t.Fatalf("buildItem: %v", err)
	}
	if it.Key != "pkg.Foo" {
		t.Errorf("Key = %q, want pkg.Foo", it.Key)
	}
	if it.State != quest.TODO {
		t.Errorf("State = %v, want TODO", it.State)
	}

	var p funcPayload
	if err := it.DecodePayload(&p); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if p.Lang != "go" || p.Root != "/root" {
		t.Errorf("payload lang=%q root=%q, want go /root", p.Lang, p.Root)
	}
	// Runtime fields reins owns are zeroed at build time.
	if p.Fn.Status != "" || p.Fn.Attempt != 0 || p.Fn.CoveragePct != 0 || p.Fn.TestMtime != "" || p.Fn.FailOutput != "" {
		t.Errorf("runtime fields not stripped: %+v", p.Fn)
	}
}
