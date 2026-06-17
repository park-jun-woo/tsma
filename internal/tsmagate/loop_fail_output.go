//ff:func feature=gate type=helper control=selection
//ff:what loopFailOutput: 루프 측정의 실행 실패 텍스트를 고른다 — runner 에러가 있으면 그 메시지, 아니면 결과 출력, 둘 다 없으면 기본 문구. measureLoop의 분기 평탄화용 분리.
package tsmagate

import "github.com/park-jun-woo/tsma/internal/runner"

// loopFailOutput picks the failure feedback text for a loop measurement: the
// runner error if present, otherwise the test output, otherwise a default. It is
// split out so measureLoop stays at nesting depth ≤2.
func loopFailOutput(err error, res *runner.Result) string {
	switch {
	case err != nil:
		return err.Error()
	case res != nil:
		return res.Output
	default:
		return "test runner returned no result"
	}
}
