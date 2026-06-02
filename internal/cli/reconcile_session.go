//ff:func feature=cli type=helper control=sequence
//ff:what 세션을 현재 소스와 동기화: 구조 머지 후 (옵션)신규 함수만 배치 측정
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// reconcileSession syncs a loaded session with the current source: it re-indexes
// and merges the function set (reconcileStructure), then — when measure is true —
// batch-measures only the newly-added functions so they get an accurate
// PASS/measured-TODO status instead of a blind TODO. It recomputes the summary
// and returns the number of functions added and removed. The session is not
// persisted here; callers decide when to Save.
func reconcileSession(root string, sess *model.Session, measure bool) (added, removed int, err error) {
	newFuncs, removed, err := reconcileStructure(root, sess)
	if err != nil {
		return 0, 0, err
	}
	if measure {
		measureFuncs(root, sess.Lang, newFuncs)
	}
	sess.RecalcSummary()
	return len(newFuncs), removed, nil
}
