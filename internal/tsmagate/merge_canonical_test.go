//ff:func feature=gate type=test
//ff:what mergeCanonical 단위테스트(BUG-002): 다함수 누적(F1 후 F2 → 둘 다 보존), 같은함수 재시도(그 함수 블록만 교체·타함수 불변), import 라인 dedup, 마커 유실 폴백(기존 전체 보존 + 추가, 소실 없음), 단일함수 회귀(누적=교체 동일 결과), 빈 정명(신규 파일 생성)을 임시 입력으로 덮는다. 언어중립 엔진이므로 go/typescript/python 토큰 분기도 함께 확인한다.

package tsmagate

import (
	"strings"
	"testing"
)

func TestMergeCanonical_MultiFuncAccumulates(t *testing.T) {
	f1 := "package pkg\n\nimport \"testing\"\n\nfunc TestF1(t *testing.T) { _ = 1 }\n"
	merged1 := mergeCanonical("", f1, "pkg.F1", "go")
	if !strings.Contains(merged1, "TestF1") {
		t.Fatalf("F1 missing after first merge:\n%s", merged1)
	}

	f2 := "package pkg\n\nimport \"testing\"\n\nfunc TestF2(t *testing.T) { _ = 2 }\n"
	merged2 := mergeCanonical(merged1, f2, "pkg.F2", "go")
	if !strings.Contains(merged2, "TestF1") {
		t.Errorf("F1 lost after F2 merge (BUG-002 regression):\n%s", merged2)
	}
	if !strings.Contains(merged2, "TestF2") {
		t.Errorf("F2 missing after merge:\n%s", merged2)
	}
	if !strings.Contains(merged2, "tsma:begin fn=pkg.F1") || !strings.Contains(merged2, "tsma:begin fn=pkg.F2") {
		t.Errorf("both function markers expected:\n%s", merged2)
	}
}

func TestMergeCanonical_SameFuncRetryReplacesOnlyItsBlock(t *testing.T) {
	f1 := "package pkg\n\nimport \"testing\"\n\nfunc TestF1(t *testing.T) { _ = 1 }\n"
	f2 := "package pkg\n\nimport \"testing\"\n\nfunc TestF2(t *testing.T) { _ = 2 }\n"
	merged := mergeCanonical(mergeCanonical("", f1, "pkg.F1", "go"), f2, "pkg.F2", "go")

	// Retry F1 with a new body — only F1's block should change.
	f1v2 := "package pkg\n\nimport \"testing\"\n\nfunc TestF1(t *testing.T) { _ = 999 }\n"
	retried := mergeCanonical(merged, f1v2, "pkg.F1", "go")

	if !strings.Contains(retried, "_ = 999") {
		t.Errorf("F1 retry body not applied:\n%s", retried)
	}
	if strings.Contains(retried, "_ = 1") {
		t.Errorf("old F1 body still present after retry:\n%s", retried)
	}
	if !strings.Contains(retried, "TestF2") || !strings.Contains(retried, "_ = 2") {
		t.Errorf("F2 block must stay unchanged after F1 retry:\n%s", retried)
	}
	if strings.Count(retried, "tsma:begin fn=pkg.F1") != 1 {
		t.Errorf("F1 block must appear exactly once (replace, not append):\n%s", retried)
	}
}

func TestMergeCanonical_ImportDedup(t *testing.T) {
	f1 := "import { describe, it } from 'vitest'\nimport { F1 } from './mod'\n\ndescribe('F1', () => { it('x', () => {}) })\n"
	f2 := "import { describe, it } from 'vitest'\nimport { F2 } from './mod'\n\ndescribe('F2', () => { it('y', () => {}) })\n"
	merged := mergeCanonical(mergeCanonical("", f1, "mod.F1", "typescript"), f2, "mod.F2", "typescript")

	if got := strings.Count(merged, "import { describe, it } from 'vitest'"); got != 1 {
		t.Errorf("shared vitest import should be deduped to 1, got %d:\n%s", got, merged)
	}
	if !strings.Contains(merged, "import { F1 } from './mod'") || !strings.Contains(merged, "import { F2 } from './mod'") {
		t.Errorf("both distinct imports must survive:\n%s", merged)
	}
	if !strings.Contains(merged, "describe('F1'") || !strings.Contains(merged, "describe('F2'") {
		t.Errorf("both describe blocks must survive:\n%s", merged)
	}
}

func TestMergeCanonical_LostMarkerFallbackPreservesAll(t *testing.T) {
	// Existing canonical with NO tsma markers (e.g. a hand-written or legacy file).
	existing := "package pkg\n\nimport \"testing\"\n\nfunc TestLegacy(t *testing.T) { _ = 7 }\n"
	newSrc := "package pkg\n\nimport \"testing\"\n\nfunc TestNew(t *testing.T) { _ = 8 }\n"
	merged := mergeCanonical(existing, newSrc, "pkg.New", "go")

	if !strings.Contains(merged, "TestLegacy") {
		t.Errorf("legacy (unmarked) content must be preserved (loss < duplication):\n%s", merged)
	}
	if !strings.Contains(merged, "TestNew") {
		t.Errorf("new block must be added:\n%s", merged)
	}
}

func TestMergeCanonical_SingleFuncRegression(t *testing.T) {
	src := "package pkg\n\nimport \"testing\"\n\nfunc TestOnly(t *testing.T) { _ = 1 }\n"
	merged := mergeCanonical("", src, "pkg.Only", "go")
	if !strings.Contains(merged, "package pkg") {
		t.Errorf("package header missing:\n%s", merged)
	}
	if !strings.Contains(merged, "TestOnly") {
		t.Errorf("single-func body missing:\n%s", merged)
	}
	// A re-merge of the same function = replace, identical block count.
	again := mergeCanonical(merged, src, "pkg.Only", "go")
	if strings.Count(again, "tsma:begin fn=pkg.Only") != 1 {
		t.Errorf("single-func re-merge must stay one block:\n%s", again)
	}
}

func TestMergeCanonical_PythonHashMarkers(t *testing.T) {
	f1 := "import pytest\n\ndef test_f1():\n    assert True\n"
	f2 := "import pytest\n\ndef test_f2():\n    assert True\n"
	merged := mergeCanonical(mergeCanonical("", f1, "mod.f1", "python"), f2, "mod.f2", "python")

	if !strings.Contains(merged, "# tsma:begin fn=mod.f1") {
		t.Errorf("python marker should use # comment:\n%s", merged)
	}
	if !strings.Contains(merged, "def test_f1") || !strings.Contains(merged, "def test_f2") {
		t.Errorf("both python tests must survive:\n%s", merged)
	}
	if got := strings.Count(merged, "import pytest"); got != 1 {
		t.Errorf("shared pytest import should dedup to 1, got %d:\n%s", got, merged)
	}
}

func TestCommentTokensFor_AllLanguages(t *testing.T) {
	cases := map[string]string{
		"go":         "//",
		"typescript": "//",
		"java":       "//",
		"csharp":     "//",
		"rust":       "//",
		"python":     "#",
		"unknown":    "//", // default → go tokens
	}
	for lang, wantLine := range cases {
		if got := commentTokensFor(lang).line; got != wantLine {
			t.Errorf("commentTokensFor(%q).line = %q, want %q", lang, got, wantLine)
		}
		if len(commentTokensFor(lang).headerPrefixes) == 0 {
			t.Errorf("commentTokensFor(%q) must carry header prefixes", lang)
		}
	}
}

func TestAssemble_EmptyHeaderAndEmptyBody(t *testing.T) {
	// Empty header → body only.
	if got := assemble(nil, "body line\n"); got != "body line\n" {
		t.Errorf("empty header assemble = %q", got)
	}
	// Empty body → header only.
	if got := assemble([]string{"package pkg"}, ""); got != "package pkg\n" {
		t.Errorf("empty body assemble = %q", got)
	}
	// Both present → blank-line separated.
	if got := assemble([]string{"package pkg"}, "func X(){}\n"); got != "package pkg\n\nfunc X(){}\n" {
		t.Errorf("both assemble = %q", got)
	}
}

func TestMergeCanonical_RustNotAccumulatedHere(t *testing.T) {
	// Rust uses // markers if ever routed here, but in practice Rust is excluded
	// at the caller. This just confirms the engine is language-neutral for rust
	// tokens without panicking and wraps with // markers.
	src := "use super::*;\n\n#[test]\nfn t_one() { assert!(true); }\n"
	merged := mergeCanonical("", src, "mod::one", "rust")
	if !strings.Contains(merged, "// tsma:begin fn=mod::one") {
		t.Errorf("rust tokens should use // marker:\n%s", merged)
	}
}

func TestReplaceOrAddBlock_AppendIntoEmptyBody(t *testing.T) {
	tok := commentTokensFor("go")
	block := wrapBlock("func TestX(t *testing.T){}", "pkg.X", tok)
	// Existing body empty → newBlock returned as-is (the base == "" branch).
	if got := replaceOrAddBlock("", block, "pkg.X", tok); got != block {
		t.Errorf("append into empty body = %q, want %q", got, block)
	}
	// Existing body present, different fn → appended with a separator.
	prior := wrapBlock("func TestY(t *testing.T){}", "pkg.Y", tok)
	merged := replaceOrAddBlock(prior, block, "pkg.X", tok)
	if !strings.Contains(merged, "fn=pkg.Y") || !strings.Contains(merged, "fn=pkg.X") {
		t.Errorf("both blocks expected:\n%s", merged)
	}
}

func TestReplaceOrAddBlock_ReplacesExistingBlock(t *testing.T) {
	tok := commentTokensFor("go")
	oldBlock := wrapBlock("func TestX(t *testing.T){ _ = 1 }", "pkg.X", tok)
	other := wrapBlock("func TestY(t *testing.T){}", "pkg.Y", tok)
	body := oldBlock + "\n\n" + other
	newBlock := wrapBlock("func TestX(t *testing.T){ _ = 2 }", "pkg.X", tok)
	got := replaceOrAddBlock(body, newBlock, "pkg.X", tok)
	if strings.Contains(got, "_ = 1") {
		t.Errorf("old block body must be replaced:\n%s", got)
	}
	if !strings.Contains(got, "_ = 2") || !strings.Contains(got, "fn=pkg.Y") {
		t.Errorf("new block and the other fn's block must both be present:\n%s", got)
	}
	if strings.Count(got, "tsma:begin fn=pkg.X") != 1 {
		t.Errorf("pkg.X block must appear exactly once (replace, not append):\n%s", got)
	}
}

func TestSplitHeader_EndMarkerBoundary(t *testing.T) {
	tok := commentTokensFor("go")
	// A file whose first content line is an end-marker (degenerate) must still
	// treat it as body boundary, not header.
	src := "package pkg\n// tsma:end fn=pkg.X\nfunc Z(){}\n"
	header, body := splitHeader(src, tok)
	if strings.Contains(strings.Join(header, "\n"), "tsma:end") {
		t.Errorf("end marker must not be absorbed into header: %v", header)
	}
	if !strings.Contains(body, "tsma:end fn=pkg.X") {
		t.Errorf("end marker must stay in body:\n%s", body)
	}
}

func TestSplitHeader_CommentLineInHeader(t *testing.T) {
	tok := commentTokensFor("go")
	// A non-marker comment line in the header region is kept in the header.
	src := "package pkg\n// a leading note\nimport \"testing\"\n\nfunc TestZ(t *testing.T){}\n"
	header, body := splitHeader(src, tok)
	hj := strings.Join(header, "\n")
	if !strings.Contains(hj, "// a leading note") {
		t.Errorf("leading comment must stay in header:\n%s", hj)
	}
	if !strings.Contains(body, "func TestZ") {
		t.Errorf("body must hold the test func:\n%s", body)
	}
}

func TestMergeCanonical_EmptyExistingCreatesFile(t *testing.T) {
	src := "package pkg\n\nfunc TestX(t *testing.T) {}\n"
	merged := mergeCanonical("", src, "pkg.X", "go")
	if !strings.HasSuffix(merged, "\n") {
		t.Errorf("merged file must end with a single newline:\n%q", merged)
	}
	if !strings.Contains(merged, "tsma:begin fn=pkg.X") || !strings.Contains(merged, "tsma:end fn=pkg.X") {
		t.Errorf("new file must carry begin+end markers:\n%s", merged)
	}
}
