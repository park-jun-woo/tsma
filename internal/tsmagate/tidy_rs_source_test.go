//ff:func feature=gate type=test lang=rust
//ff:what tidyRsSource 단위테스트: rustfmt 부재(LookPath 실패 → 입력 그대로),
// 성공(가짜 rustfmt가 stdin을 소비하고 포맷 결과를 내면 그 출력을 반환),
// 실패(비-0 종료 → 입력 그대로), 공백 출력(→ 입력 그대로)의 전 분기를 PATH에
// 심은 가짜 rustfmt 스크립트로 결정적으로 덮는다 — 실제 rustfmt 불필요.
package tsmagate

import "testing"

// TestTidyRsSourceFormatterAbsent covers the LookPath-failure branch: without
// rustfmt on PATH the input is returned unchanged.
func TestTidyRsSourceFormatterAbsent(t *testing.T) {
	emptyPath(t)
	src := "fn x() {}\n"
	if got := tidyRsSource(src); got != src {
		t.Errorf("without rustfmt the source must be unchanged: %q", got)
	}
}

// TestTidyRsSourceFormatSuccess drives the happy path with a fake rustfmt that
// consumes stdin and emits fixed formatted output.
func TestTidyRsSourceFormatSuccess(t *testing.T) {
	installFakeTool(t, "rustfmt", "#!/bin/sh\ncat > /dev/null\nprintf 'fn fmted() {}\\n'\n")
	if got := tidyRsSource("fn x() {}\n"); got != "fn fmted() {}\n" {
		t.Errorf("tidyRsSource = %q, want %q", got, "fn fmted() {}\n")
	}
}

// TestTidyRsSourceRunError covers the cmd.Run failure branch: a fake rustfmt
// that exits non-zero degrades to the input unchanged.
func TestTidyRsSourceRunError(t *testing.T) {
	installFakeTool(t, "rustfmt", "#!/bin/sh\ncat > /dev/null\nexit 2\n")
	src := "fn x() {}\n"
	if got := tidyRsSource(src); got != src {
		t.Errorf("run failure must return src unchanged, got %q", got)
	}
}

// TestTidyRsSourceEmptyOutput covers the blank-output half of the fallback
// condition: a fake rustfmt that succeeds but prints only whitespace degrades to
// the input unchanged.
func TestTidyRsSourceEmptyOutput(t *testing.T) {
	installFakeTool(t, "rustfmt", "#!/bin/sh\ncat > /dev/null\nprintf '  '\n")
	src := "fn x() {}\n"
	if got := tidyRsSource(src); got != src {
		t.Errorf("blank rustfmt output must return src unchanged, got %q", got)
	}
}
