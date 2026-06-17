//ff:func feature=gate type=helper control=sequence level=error
//ff:what testTargetPath: 생성 소스를 쓸 테스트 파일을 §2-2 우선순위로 정해 root-상대 경로로 돌려준다. ①함수에 이미 귀속된 테스트(tm.Files) → 첫 매칭 파일 덮어쓰기, ②디스크의 오명명 변형(FindMisnamedTest) → 정명으로 rename해 흡수, ③테스트 전무한 신규 함수 → match.CanonicalTestPath로 정명 도출. rename 실패·도출 불가는 error로 돌려 Prepare가 TestFailed로 드러낸다(무음 금지).

package tsmagate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/match"
)

// testTargetPath decides which test file the generated source is written to,
// following the §2-2 priority and returning a project-root-relative path:
//
//  1. Tests already attributed to the function (tm.Files non-empty) → write to
//     the first matched file (overwrite — the LLM regenerates the whole file).
//  2. A misnamed variant exists on disk (match.FindMisnamedTest found) → absorb
//     it by renaming to the canonical name, then write there (so a stale
//     misnamed file is not left behind to double-match).
//  3. A brand-new function with no test of any kind → derive the canonical path
//     directly from the source file via match.CanonicalTestPath.
//
// found/tm come from the caller's pre-write MatchFunc so the decision and the
// later re-match see a consistent view. A rename or an underivable path is
// returned as an error so Prepare can surface it as TestFailed (never silent).
func testTargetPath(p funcPayload, tm match.TestMatch, found bool) (string, error) {
	// 1. Already-attributed test file → overwrite it.
	if found && len(tm.Files) > 0 {
		return tm.Files[0], nil
	}
	// 2. Misnamed variant on disk → absorb into the canonical name.
	if misnamed, canonical, ok := match.FindMisnamedTest(p.Root, p.Lang, p.Fn.File); ok {
		absMis := filepath.Join(p.Root, misnamed)
		absCanon := filepath.Join(p.Root, canonical)
		if err := os.Rename(absMis, absCanon); err != nil {
			return "", fmt.Errorf("absorb misnamed test %s → %s: %w", misnamed, canonical, err)
		}
		return canonical, nil
	}
	// 3. Brand-new function → derive the canonical path from the source file.
	canonical := match.CanonicalTestPath(p.Lang, p.Fn.File)
	if canonical == "" {
		return "", fmt.Errorf("cannot derive a test path for %s (lang=%s)", p.Fn.File, p.Lang)
	}
	return canonical, nil
}
