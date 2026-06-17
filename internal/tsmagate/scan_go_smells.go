//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what scanGoSmells: Prepare에서 분리한 Go 테스트 smell 스캔. 매칭된 root-상대 테스트 파일들을 root와 join해 smell.ScanGo로 정적 스캔하고 Finding을 합쳐 돌려준다. 파싱 에러는 무시(continue)한다 — 깨진 테스트 파일은 tests-must-pass가 잡는다. Prepare의 if/for/if 3중 중첩을 평탄화하려 추출(Q1 depth ≤2).

package tsmagate

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/smell"
)

// scanGoSmells statically scans the matched Go test files for escape-hatch
// smells and returns the collected findings. tm.Files are root-relative, so they
// are joined with root. A parse error is ignored (continue): a broken test file
// is judged by tests-must-pass, not here. It is extracted from Prepare purely to
// flatten the lang-guard/file-loop/parse-guard nesting.
func scanGoSmells(root string, files []string) []smell.Finding {
	var findings []smell.Finding
	for _, file := range files {
		found, err := smell.ScanGo(filepath.Join(root, file))
		if err != nil {
			continue
		}
		findings = append(findings, found...)
	}
	return findings
}
