//ff:func feature=gate type=test
//ff:what smell 스캐너 append 분기·scanSmells java/csharp 디스패치 단위테스트.
// 실제 tree-sitter CLI 대신 TSMA_TREE_SITTER로 가짜 스크립트(고정 최소 XML 출력)를
// 주입해, ScanTS/ScanRs/ScanJava/ScanCs가 err=nil로 돌아와 scan*Smells의
// findings-append 문이 결정적으로 실행되게 한다(빈 program 트리 → 발견 0건).
// scanSmells의 java/csharp 아암도 같은 경로로 함께 덮는다.
package tsmagate

import (
	"os"
	"path/filepath"
	"testing"
)

// installFakeTreeSitter installs a fake tree-sitter CLI that emits a minimal
// valid `parse --xml` document (one source, an empty program root) and exports
// it via TSMA_TREE_SITTER, so every ScanX succeeds with zero findings and the
// scanners' append statement runs without the real CLI or grammars.
func installFakeTreeSitter(t *testing.T) {
	t.Helper()
	script := `#!/bin/sh
cat <<'EOF'
<?xml version="1.0"?>
<sources>
  <source name="/tsma-fake/only.src">
    <program srow="0" scol="0" erow="0" ecol="0"></program>
  </source>
</sources>
EOF
`
	dir := t.TempDir()
	path := filepath.Join(dir, "tree-sitter")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TSMA_TREE_SITTER", path)
}

// TestScanTSSmells_AppendWithFakeCLI covers the findings-append statement: the
// fake CLI parses successfully (empty tree), so the scan returns without error
// and appends its (empty) findings.
func TestScanTSSmells_AppendWithFakeCLI(t *testing.T) {
	installFakeTreeSitter(t)
	if got := scanTSSmells(t.TempDir(), []string{"a.test.ts"}); len(got) != 0 {
		t.Errorf("empty parse tree must yield no findings, got %+v", got)
	}
}

// TestScanRsSmells_AppendWithFakeCLI covers the findings-append statement for
// the Rust scanner via the same fake CLI.
func TestScanRsSmells_AppendWithFakeCLI(t *testing.T) {
	installFakeTreeSitter(t)
	if got := scanRsSmells(t.TempDir(), []string{"lib.rs"}); len(got) != 0 {
		t.Errorf("empty parse tree must yield no findings, got %+v", got)
	}
}

// TestScanJavaSmells_AppendWithFakeCLI covers the findings-append statement for
// the Java scanner via the same fake CLI.
func TestScanJavaSmells_AppendWithFakeCLI(t *testing.T) {
	installFakeTreeSitter(t)
	if got := scanJavaSmells(t.TempDir(), []string{"FooTest.java"}); len(got) != 0 {
		t.Errorf("empty parse tree must yield no findings, got %+v", got)
	}
}

// TestScanCsSmells_AppendWithFakeCLI covers the findings-append statement for
// the C# scanner via the same fake CLI.
func TestScanCsSmells_AppendWithFakeCLI(t *testing.T) {
	installFakeTreeSitter(t)
	if got := scanCsSmells(t.TempDir(), []string{"FooTests.cs"}); len(got) != 0 {
		t.Errorf("empty parse tree must yield no findings, got %+v", got)
	}
}

// TestScanSmells_JavaAndCsharpDispatch covers the java and csharp arms of the
// scanSmells language switch (the go/typescript/rust arms are covered
// elsewhere), routing through the fake CLI so the dispatched scanners succeed.
func TestScanSmells_JavaAndCsharpDispatch(t *testing.T) {
	installFakeTreeSitter(t)
	if got := scanSmells("java", t.TempDir(), []string{"FooTest.java"}); len(got) != 0 {
		t.Errorf("scanSmells(java) = %+v, want none", got)
	}
	if got := scanSmells("csharp", t.TempDir(), []string{"FooTests.cs"}); len(got) != 0 {
		t.Errorf("scanSmells(csharp) = %+v, want none", got)
	}
}
