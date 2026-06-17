//ff:func feature=gate type=helper control=sequence lang=go
//ff:what prepareLoopGo: Go 루프 분기의 비침습 측정 파이프라인(C2/C3). 생성 raw를 sanitize+tidy→go/parser 파싱(실패=잘림 C3: TestFailed+Truncated, 파일 미기록)→AST에서 테스트함수 직접 추출→잘못된 테스트명(go test가 무시할 TestXxx 위반)이면 측정 전 reject(C1 보강)→backing(.tsma/test) 기록 + 가상경로 overlay TestMatch 직접 구성(MatchFunc 우회)→smell은 backing 직접 스캔→runner/checker가 -overlay -vet=off로 소스 트리 무침습 측정→종결 처리(materialize 후/실패 후 backing·overlay JSON 정리, C2)는 finalizeBacking. measurement를 돌려준다.
package tsmagate

import "github.com/park-jun-woo/reins/pkg/quest"

// prepareLoopGo runs the Go loop's non-invasive measurement pipeline (C2/C3). It
// sanitizes+tidies the generated source, parses it (a parse failure is a
// truncated output → C3: TestFailed+Truncated, nothing written), extracts the
// test functions straight from the AST (no disk re-match), writes a backing file
// and hand-builds the overlay TestMatch, scans the backing for smells, measures
// via the Go runner/checker with -overlay -vet=off (the source tree is never
// touched), and only at a terminal pass materializes the backing to the canonical
// path. It returns the measurement for the gate rules.
func prepareLoopGo(it *quest.Item, p funcPayload, raw []byte) *measurement {
	m := &measurement{FuncName: p.Fn.QualifiedName}
	src := sanitizeGoSource(string(raw))
	funcs, ok := parseTestFuncs(src)
	if !ok {
		m.TestFailed, m.Truncated = true, true
		return m
	}
	if bad, malformed := firstMalformedTestName(funcs); malformed {
		m.TestFailed = true
		m.FailExpected = "well-formed Go test names (TestXxx — uppercase after 'Test')"
		m.FailOutput = "test name `" + bad + "` is malformed: the character after 'Test' must be uppercase, " +
			"else `go test` silently skips it and nothing runs"
		return m
	}
	tm, backingRel, err := buildLoopTestMatch(p, it, src, funcs)
	if err != nil {
		m.TestFailed, m.FailOutput = true, err.Error()
		return m
	}
	m.Smells = scanGoSmells(p.Root, []string{backingRel})
	measureLoop(m, p, tm)
	finalizeBacking(p, it, m, backingRel)
	return m
}
