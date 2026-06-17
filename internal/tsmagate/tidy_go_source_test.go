//ff:func feature=gate type=test
//ff:what tidy_go_source 단위테스트: tidyGoSource가 미사용 import 제거(+컴파일 가능), 모든 import 사용 시 의미 보존(gofmt만), 잘린/파싱불가 소스 패스스루(패닉 없음), blank/dot import 보존, 그리고 sanitize+tidy 합성(펜스+미사용import → 순수·정리된 Go)을 검증한다.

package tsmagate

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// mustCompileParse asserts the source parses (a proxy for "compilable shape":
// after pruning, leftover unused imports would still parse, so we additionally
// assert the dropped package is absent).
func mustParse(t *testing.T, src string) {
	t.Helper()
	if _, err := parser.ParseFile(token.NewFileSet(), "", src, parser.AllErrors); err != nil {
		t.Fatalf("result does not parse: %v\n%s", err, src)
	}
}

func TestTidyGoSource_RemovesUnusedImport(t *testing.T) {
	src := "package pkg\n\nimport (\n\t\"fmt\"\n\t\"strings\"\n)\n\nfunc Use() { fmt.Println(\"hi\") }\n"
	got := tidyGoSource(src)
	mustParse(t, got)
	if strings.Contains(got, "strings") {
		t.Fatalf("unused \"strings\" import should be removed:\n%s", got)
	}
	if !strings.Contains(got, "\"fmt\"") {
		t.Fatalf("used \"fmt\" import must survive:\n%s", got)
	}
}

func TestTidyGoSource_DropsEmptyImportBlock(t *testing.T) {
	// The only import is unused, so the whole import block must vanish.
	src := "package pkg\n\nimport \"strings\"\n\nfunc F() {}\n"
	got := tidyGoSource(src)
	mustParse(t, got)
	if strings.Contains(got, "import") {
		t.Fatalf("empty import block should be removed entirely:\n%s", got)
	}
}

func TestTidyGoSource_PreservesAllUsed(t *testing.T) {
	src := "package pkg\n\nimport (\n\t\"fmt\"\n\t\"strings\"\n)\n\nfunc F() string { return fmt.Sprint(strings.ToUpper(\"x\")) }\n"
	got := tidyGoSource(src)
	mustParse(t, got)
	if !strings.Contains(got, "\"fmt\"") || !strings.Contains(got, "\"strings\"") {
		t.Fatalf("both used imports must be kept:\n%s", got)
	}
}

func TestTidyGoSource_TruncatedPassthrough(t *testing.T) {
	// Unparsable (truncated) source: returned unchanged, no panic.
	src := "package pkg\n\nimport (\n\t\"fmt\"\n\nfunc F() { fmt.Println("
	got := tidyGoSource(src)
	if got != src {
		t.Fatalf("truncated source must be returned unchanged:\ngot  %q\nwant %q", got, src)
	}
}

func TestTidyGoSource_PreservesBlankAndDotImports(t *testing.T) {
	src := "package pkg\n\nimport (\n\t_ \"embed\"\n\t. \"strings\"\n\t\"bytes\"\n)\n\nfunc F() { _ = ToUpper(\"x\") }\n"
	got := tidyGoSource(src)
	mustParse(t, got)
	if !strings.Contains(got, "_ \"embed\"") {
		t.Fatalf("blank import must be preserved:\n%s", got)
	}
	if !strings.Contains(got, ". \"strings\"") {
		t.Fatalf("dot import must be preserved:\n%s", got)
	}
	if strings.Contains(got, "\"bytes\"") {
		t.Fatalf("unused \"bytes\" must be removed even alongside blank/dot:\n%s", got)
	}
}

func TestTidyGoSource_PreservesAliasUsed(t *testing.T) {
	src := "package pkg\n\nimport f \"fmt\"\n\nfunc F() { f.Println(\"x\") }\n"
	got := tidyGoSource(src)
	mustParse(t, got)
	if !strings.Contains(got, "f \"fmt\"") {
		t.Fatalf("used aliased import must be preserved:\n%s", got)
	}
}

func TestSanitizeGoSource_FenceAndUnusedImportCombined(t *testing.T) {
	// sanitize unwraps the fence, tidy prunes the unused "strings": the disk
	// content must be pure, tidy Go with no fence and no dead import.
	in := "Here you go:\n```go\npackage pkg\n\nimport (\n\t\"fmt\"\n\t\"strings\"\n)\n\nfunc TestX(t *testing.T) { fmt.Println(\"ok\") }\n```\nThanks!"
	got := sanitizeGoSource(in)
	mustParse(t, got)
	if strings.Contains(got, "```") || strings.Contains(got, "Here you go") {
		t.Fatalf("fence/prose must be stripped:\n%s", got)
	}
	if strings.Contains(got, "strings") {
		t.Fatalf("unused import must be pruned:\n%s", got)
	}
	if !strings.HasPrefix(got, "package pkg") || !strings.HasSuffix(got, "\n") {
		t.Fatalf("result should be a clean Go file:\n%s", got)
	}
}
