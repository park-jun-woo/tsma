//ff:func feature=gate type=helper control=sequence
//ff:what readFile: absPath의 파일 내용을 돌려준다(에러면 ""). read-only — Render가 기존 테스트 파일 전체를 프롬프트에 실어 LLM이 재생성 시 형제 테스트를 무음으로 떨구지 않게 한다.

package tsmagate

import "os"

// readFile returns the file content at absPath, or "" on error. Read-only.
func readFile(absPath string) string {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return ""
	}
	return string(data)
}
