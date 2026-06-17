//ff:func feature=gate type=helper control=sequence level=error
//ff:what Render: 한 함수에 대해 "테스트 작성 + 브랜치 100% 달성" 작성 프롬프트를 문자열로 반환(read-only, s.Meta 변경 금지). 구버전 print_todo_function/print_rename_instruction의 텍스트를 출력이 아니라 return 문자열로 재표현한다. 매칭된 테스트 파일·오명명 rename 힌트를 안내(파일 읽기만, 부작용 없음).

package tsmagate

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/match"
)

// Render returns the authoring prompt `next` shows for one function. It is
// read-only: it reads the source/test files to locate the matched test (and an
// optional rename hint) but never runs tests, never mutates s.Meta, and never
// writes disk. The text re-expresses the legacy print_todo_function /
// print_rename_instruction output as a returned string instead of stdout.
func (d *Definition) Render(s *quest.Session, it *quest.Item) (string, error) {
	var p funcPayload
	if err := it.DecodePayload(&p); err != nil {
		return "", fmt.Errorf("decode payload for %s: %w", it.Key, err)
	}
	fn := p.Fn

	var b strings.Builder
	fmt.Fprintf(&b, "%s  TODO\n", fn.QualifiedName)
	fmt.Fprintf(&b, "  file: %s:%d-%d\n", fn.File, fn.StartLine, fn.EndLine)

	// Locate the test file (content-aware for Go, file-name based otherwise).
	tm, found := match.NewFuncMatcher(p.Lang).MatchFunc(p.Root, &fn)
	if found && len(tm.Files) > 0 {
		fmt.Fprintf(&b, "  test: %s\n", strings.Join(tm.Files, ", "))
	} else {
		fmt.Fprintf(&b, "  test: (not found)\n")
		// Surface a rename hint when a near-miss test file exists on disk.
		if misnamed, canonical, ok := match.FindMisnamedTest(p.Root, p.Lang, fn.File); ok {
			fmt.Fprintf(&b, "  ▶ Test file misnamed: rename `%s` → `%s`.\n", misnamed, canonical)
		}
	}

	fmt.Fprintf(&b, "\nWrite tests for %s that achieve 100%% branch coverage.\n", fn.Name)
	b.WriteString("Cover every branch (if/else, switch, error paths, loops).\n")
	b.WriteString("Then submit so the gate re-measures branch coverage on disk.\n")
	return b.String(), nil
}
