//ff:func feature=gate type=helper control=sequence level=error
//ff:what writeTestFile: src를 root 아래 testRel(root-상대)에 쓰고, 부모 디렉터리는 필요시 생성한다. 기존 파일은 덮어쓴다 — 루프 모델상 LLM이 매 시도 완전한 _test.go를 내므로 최신 생성이 이긴다. 에러는 전파해 Prepare가 TestFailed 측정으로 드러낸다(쓰기 부작용을 무음으로 삼키지 않는다).

package tsmagate

import (
	"fmt"
	"os"
	"path/filepath"
)

// writeTestFile writes src to testRel (project-root-relative) under root,
// creating the parent directory if needed. It overwrites any existing file —
// the loop's model is that the LLM emits a complete _test.go each attempt, so
// the latest generation wins. Errors propagate so Prepare surfaces them as a
// TestFailed measurement (the write side effect is never silently swallowed).
func writeTestFile(root, testRel, src string) error {
	abs := filepath.Join(root, testRel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return fmt.Errorf("create test dir for %s: %w", testRel, err)
	}
	if err := os.WriteFile(abs, []byte(src), 0o644); err != nil {
		return fmt.Errorf("write test file %s: %w", testRel, err)
	}
	return nil
}
