//ff:func feature=gate type=helper control=sequence level=error
//ff:what Render: 한 함수에 대한 생성형 테스트-작성 프롬프트를 문자열로 반환(read-only, s.Meta 변경 금지). 언어·함수 풀네임/시그니처·패키지명·소스 본문(file:line)·기존 테스트(있으면)를 1회용 컨텍스트로 실어 "완전한 _test.go를 작성하라"고 지시한다. loop가 이 출력을 backend.Complete의 user 프롬프트로 쓴다(미커버 file:line은 reins가 FAIL 피드백으로 따로 덧붙이므로 중복 주지 않음). 매칭된 테스트 파일·오명명 rename 힌트도 안내(파일 읽기만, 부작용 없음).

package tsmagate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/match"
)

// Render returns the generative authoring prompt for one function. It is
// read-only: it reads the source/test files to assemble one-shot context (the
// function's body, its package, and any existing test) but never runs tests,
// never mutates s.Meta, and never writes disk. The loop feeds this string to the
// LLM backend as the user prompt; on retry, reins appends the uncovered file:line
// feedback (renderVerdictText) itself, so Render deliberately does not enumerate
// uncovered lines (no double exposure).
func (d *Definition) Render(s *quest.Session, it *quest.Item) (string, error) {
	var p funcPayload
	if err := it.DecodePayload(&p); err != nil {
		return "", fmt.Errorf("decode payload for %s: %w", it.Key, err)
	}
	fn := p.Fn

	var b strings.Builder
	fmt.Fprintf(&b, "Language: %s\n", p.Lang)
	fmt.Fprintf(&b, "Function: %s\n", fn.QualifiedName)
	if fn.Receiver != "" {
		fmt.Fprintf(&b, "Receiver: %s\n", fn.Receiver)
	}
	fmt.Fprintf(&b, "Source: %s:%d-%d\n", fn.File, fn.StartLine, fn.EndLine)

	// Package declaration (Go white-box default: the generated test shares the
	// source package so it can reach unexported symbols). Best-effort: a read
	// failure simply omits the hint rather than failing the whole render.
	srcAbs := filepath.Join(p.Root, fn.File)
	if pkg := goPackageName(srcAbs); pkg != "" {
		fmt.Fprintf(&b, "Package: %s\n", pkg)
	}

	// The function body itself — the primary context the LLM needs to enumerate
	// branches. Pulled by line range so unrelated code in the file is excluded.
	if body := readLineRange(srcAbs, fn.StartLine, fn.EndLine); body != "" {
		b.WriteString("\nFunction under test:\n")
		fmt.Fprintf(&b, "```%s\n%s\n```\n", p.Lang, body)
	}

	// Locate the test file (content-aware for Go, file-name based otherwise) and,
	// when one exists, include its current content so the LLM extends rather than
	// silently drops sibling tests when it regenerates the file.
	tm, found := match.NewFuncMatcher(p.Lang).MatchFunc(p.Root, &fn)
	if found && len(tm.Files) > 0 {
		fmt.Fprintf(&b, "\nExisting test file: %s\n", strings.Join(tm.Files, ", "))
		if existing := readFile(filepath.Join(p.Root, tm.Files[0])); existing != "" {
			fmt.Fprintf(&b, "```%s\n%s\n```\n", p.Lang, existing)
		}
	} else {
		b.WriteString("\nExisting test file: (none)\n")
		// Surface a rename hint when a near-miss test file exists on disk.
		if misnamed, canonical, ok := match.FindMisnamedTest(p.Root, p.Lang, fn.File); ok {
			fmt.Fprintf(&b, "Note: a misnamed test `%s` exists; it will be absorbed into `%s`.\n", misnamed, canonical)
		}
	}

	fmt.Fprintf(&b, "\nWrite a complete, compilable test file for %s that achieves 100%% branch coverage.\n", fn.Name)
	b.WriteString("Cover every branch (if/else, switch, error paths, loops).\n")
	b.WriteString("Output ONLY the full test file source — no prose, no markdown fences.\n")
	return b.String(), nil
}
