//ff:func feature=gate type=helper control=sequence
//ff:what sanitizeGoSource: LLM이 테스트 출력에 두른 래퍼(마크다운 펜스 ```go … ```, 펜스 밖 산문)를 벗겨 디스크에 순수 Go만 남긴다. 펜스가 없으면 앞뒤 공백만 트림한다. validator가 아니라 best-effort 정리 — 컴파일 안 되는 잔여물은 하류 tests-must-pass 게이트가 FAIL로 잡아 되먹임하므로 sanitize는 거부하지 않고 풀기만 한다.

package tsmagate

import "strings"

// sanitizeGoSource strips the wrappers an LLM may add around a generated test
// file so what lands on disk is pure Go. It removes Markdown code fences
// (```go … ```) and any prose lines outside the fenced block; when no fence is
// present it returns the input trimmed of surrounding whitespace. This is a
// best-effort cleanup, not a validator: anything that still does not compile is
// caught downstream by the tests-must-pass gate (a compile failure FAILs and is
// fed back), so sanitize never needs to reject — it only unwraps.
func sanitizeGoSource(raw string) string {
	s := strings.ReplaceAll(raw, "\r\n", "\n")
	// If a fenced block exists, keep only its contents (the first fenced block).
	if i := strings.Index(s, "```"); i >= 0 {
		// Skip to the end of the opening fence line (drops the optional "go" tag).
		rest := s[i+3:]
		if nl := strings.IndexByte(rest, '\n'); nl >= 0 {
			rest = rest[nl+1:]
		}
		// Cut at the closing fence.
		if j := strings.Index(rest, "```"); j >= 0 {
			rest = rest[:j]
		}
		return strings.TrimSpace(rest) + "\n"
	}
	return strings.TrimSpace(s) + "\n"
}
