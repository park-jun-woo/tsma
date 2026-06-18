//ff:type feature=gate type=model
//ff:what commentTokens는 언어중립 누적 엔진(merge_canonical)이 함수별 마커 블록을 찍고 헤더 영역을 가르는 데 필요한 언어별 사실을 데이터로 운반한다: line=마커 라인의 라인-주석 접두사(`//` 또는 `#`), headerPrefixes=라인이 헤더(import/package/use/using…) 영역에 속하는지 판정하는 선두 키워드 집합. 코드가 아닌 데이터로 두어 엔진을 언어중립으로 유지한다. commentTokensFor가 lang→이 구조체로 매핑한다.

package tsmagate

// commentTokens carries the per-language facts the language-neutral accumulation
// engine needs: the line-comment prefix used to stamp the fn= block markers, and
// the set of leading keywords that mark a line as belonging to the header
// (import/use/using/package…) region rather than a test body. Keeping these as
// data (not code) lets merge_canonical stay language-neutral.
type commentTokens struct {
	// line is the line-comment prefix for marker lines ("//" or "#").
	line string
	// headerPrefixes are line-leading keywords (after trimming) that mark a line
	// as part of the header (import/package) region for the line-union dedup.
	headerPrefixes []string
}
