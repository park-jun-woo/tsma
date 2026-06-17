//ff:func feature=runner type=helper control=sequence level=error
//ff:what GoOverlayArgs: TestMatch.Overlay(가상경로→backing 절대경로)를 .tsma/test/overlay.json에 직렬화하고 그 오버레이를 켜는 go test 플래그(-overlay <json> -vet=off)를 돌려준다. overlay가 비면 nil을 돌려 수동 submit 경로는 현행 vet-on 그대로 둔다(스코프 분리). -vet=off 필수: 기본 vet 단계가 -overlay를 따르지 않아 가상 _test.go를 못 찾고 build failed로 죽는다(go1.22 실측).
package runner

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// GoOverlayArgs serializes the overlay map to .tsma/test/overlay.json under the
// project root and returns the `go test` flags that activate it. It returns nil
// (no error) when the overlay is empty, so the manual-submit path keeps its
// current vet-on behavior. `-vet=off` is required: the default vet step does not
// honor -overlay and fails the build looking for the virtual _test.go (go1.22).
func GoOverlayArgs(projectRoot string, overlay map[string]string) ([]string, error) {
	if len(overlay) == 0 {
		return nil, nil
	}
	jsonPath := filepath.Join(projectRoot, ".tsma", "test", "overlay.json")
	if err := os.MkdirAll(filepath.Dir(jsonPath), 0o755); err != nil {
		return nil, err
	}
	data, err := json.Marshal(struct{ Replace map[string]string }{Replace: overlay})
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(jsonPath, data, 0o644); err != nil {
		return nil, err
	}
	return []string{"-overlay", jsonPath, "-vet=off"}, nil
}
