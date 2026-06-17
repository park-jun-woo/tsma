//ff:func feature=gate type=helper control=sequence
//ff:what readLineRange: absPath 파일의 [start,end] 줄(1-based, 양끝 포함)을 돌려준다(읽기/범위 에러면 ""). read-only·best-effort — Render가 함수 본문만 줄 범위로 떼어 프롬프트 컨텍스트에 실어 파일 내 무관 코드를 배제한다.

package tsmagate

import (
	"os"
	"strings"
)

// readLineRange returns lines [start,end] (1-based, inclusive) of the file at
// absPath, or "" on any read/range error. Read-only, best-effort context.
func readLineRange(absPath string, start, end int) string {
	data, err := os.ReadFile(absPath)
	if err != nil || start < 1 || end < start {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	if start > len(lines) {
		return ""
	}
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[start-1:end], "\n")
}
