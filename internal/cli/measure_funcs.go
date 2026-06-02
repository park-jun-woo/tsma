//ff:func feature=cli type=helper control=selection
//ff:what 함수 부분집합을 언어별 배치 측정 코어로 디스패치(신규 reconcile 함수 측정용)
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// measureFuncs batch-measures coverage for an explicit subset of functions,
// dispatching Go vs non-Go. It reuses the first-scan batch cores, so a
// 100%-covered function becomes PASS and a matched partial becomes a measured
// TODO (Attempt=1); the matched test files/mtime are recorded by applyBatchReport
// so no separate attributeTests pass is needed. Used by reconcileSession to
// measure only the functions added by a rescan, leaving already-measured
// functions untouched.
func measureFuncs(root, lang string, funcs []*model.Function) {
	if len(funcs) == 0 {
		return
	}
	if lang == "go" {
		batchMeasureGoFuncs(root, lang, funcs)
		return
	}
	batchMeasureOtherFuncs(root, lang, funcs)
}
