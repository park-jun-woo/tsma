//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what goPackageName: absPath의 Go 소스에 선언된 패키지명을 돌려준다(읽기 실패·package 라인 없음이면 ""). best-effort·read-only — Render가 화이트박스 기본값(생성 테스트가 소스 패키지를 공유해 비공개 심볼 접근)을 프롬프트에 실을 때 쓴다.

package tsmagate

import (
	"os"
	"strings"
)

// goPackageName returns the package name declared in the Go source at absPath, or
// "" if it cannot be read or has no package line. Best-effort and read-only.
func goPackageName(absPath string) string {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "package ") {
			return strings.TrimSpace(strings.TrimPrefix(t, "package "))
		}
	}
	return ""
}
