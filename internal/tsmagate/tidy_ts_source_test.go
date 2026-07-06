//ff:func feature=gate type=test lang=typescript
//ff:what tidyTSSource 단위테스트: PATH에 심은 가짜 npx로 성공(포맷 출력 반환),
// 실패(비-0 종료 → 입력 그대로), 공백 출력(→ 입력 그대로) 분기를 결정적으로
// 덮는다 — 실제 prettier/npx 설치 여부와 무관하게 전 분기 도달.
package tsmagate

import "testing"

// TestTidyTSSourceFormatSuccess drives the happy path with a fake npx that
// consumes stdin and emits fixed formatted output.
func TestTidyTSSourceFormatSuccess(t *testing.T) {
	installFakeTool(t, "npx", "#!/bin/sh\ncat > /dev/null\nprintf 'const fmted = 1;\\n'\n")
	if got := tidyTSSource("const x=1;\n"); got != "const fmted = 1;\n" {
		t.Errorf("tidyTSSource = %q, want %q", got, "const fmted = 1;\n")
	}
}

// TestTidyTSSourceRunError covers the cmd.Run failure branch: a fake npx that
// exits non-zero degrades to the input unchanged.
func TestTidyTSSourceRunError(t *testing.T) {
	installFakeTool(t, "npx", "#!/bin/sh\ncat > /dev/null\nexit 1\n")
	src := "const x=1;\n"
	if got := tidyTSSource(src); got != src {
		t.Errorf("run failure must return src unchanged, got %q", got)
	}
}

// TestTidyTSSourceEmptyOutput covers the blank-output branch: a fake npx that
// succeeds but prints only whitespace degrades to the input unchanged.
func TestTidyTSSourceEmptyOutput(t *testing.T) {
	installFakeTool(t, "npx", "#!/bin/sh\ncat > /dev/null\nprintf '  '\n")
	src := "const x=1;\n"
	if got := tidyTSSource(src); got != src {
		t.Errorf("blank prettier output must return src unchanged, got %q", got)
	}
}
