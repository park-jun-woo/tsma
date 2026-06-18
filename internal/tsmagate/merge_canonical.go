//ff:func feature=gate type=helper control=sequence lang=go
//ff:what mergeCanonical: BUG-002 누적 엔진(언어중립). 기존 정명 파일 내용(없으면 빈)과 새 생성 src를, 대상 함수 QualifiedName으로 함수별 마커 블록 누적해 합친다. 동작: 새 src를 헤더/본문으로 가르고(splitHeader) 본문을 fn=<QN> 마커로 감싼 뒤(wrapBlock), 기존 본문에서 같은 fn 블록이 있으면 교체·없으면 추가(replaceOrAddBlock); 헤더는 기존+새 라인 합집합 dedup(mergeHeaders). 결과는 재포맷 없이 그대로 writeTestFile로 쓰여야 한다(머지 후 포맷 금지 — 마커 생존 불변식). Rust는 인파일 mod 누적이라 이 엔진 비대상(호출처에서 제외).

package tsmagate

import "strings"

// mergeCanonical accumulates the newly generated test src into the existing
// canonical file content, keyed by the target function's QualifiedName. It is
// language-neutral: the only per-language facts come from commentTokensFor(lang)
// (comment syntax + header keywords). The new src is split into header and body;
// the body is wrapped in fn=<qn> markers and either replaces the same function's
// prior block (same-function retry) or is appended (a different function in the
// same source file). Headers are line-unioned (dedup). The result is returned as
// already-formatted text and MUST be written verbatim (no re-format) so the
// markers survive — the core invariant of Phase 007.
//
// Rust is NOT a caller of this engine: its unit tests accumulate inside an in-file
// #[cfg(test)] mod, so BUG-002 does not apply (Phase 007 §1.1).
func mergeCanonical(existing, src, qn, lang string) string {
	tok := commentTokensFor(lang)

	newHeader, newBody := splitHeader(src, tok)
	newBlock := wrapBlock(newBody, qn, tok)

	// Brand-new canonical file: header + the single function block.
	if strings.TrimSpace(existing) == "" {
		return assemble(newHeader, newBlock)
	}

	existHeader, existBody := splitHeader(existing, tok)
	mergedHeader := mergeHeaders(existHeader, newHeader)
	mergedBody := replaceOrAddBlock(existBody, newBlock, qn, tok)
	return assemble(mergedHeader, mergedBody)
}
