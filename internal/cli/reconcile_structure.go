//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what 소스를 재인덱싱해 세션 함수 집합을 add/remove/메타갱신으로 머지(진행상태 보존)
package cli

import (
	"fmt"
	"time"

	"github.com/park-jun-woo/tsma/internal/index"
	"github.com/park-jun-woo/tsma/internal/model"
)

// reconcileStructure re-indexes the current source tree and merges the result
// into sess.Functions keyed by QualifiedName. A function present in source but
// not in the session is appended as a fresh TODO; a session function no longer
// present in source is dropped; a function present in both keeps its progress
// (Status/CoveragePct/Attempt/TestFiles/TestFile/TestMtime/FailOutput) while its
// positional metadata (File/StartLine/EndLine/Exported) is refreshed from the
// new index so coverage measurement stays accurate after edits. It returns
// pointers to the newly-added functions (re-collected after append, so they are
// valid even if the slice reallocated) and the number removed. It does NOT
// persist the session — that is the caller's responsibility.
func reconcileStructure(root string, sess *model.Session) (added []*model.Function, removed int, err error) {
	idxr := index.NewIndexer(sess.Lang)
	current, err := idxr.Index(root)
	if err != nil {
		return nil, 0, fmt.Errorf("reindex: %w", err)
	}

	cur := make(map[string]*model.Function, len(current))
	for i := range current {
		if _, dup := cur[current[i].QualifiedName]; dup {
			continue // first occurrence wins (matches FindFunction semantics)
		}
		cur[current[i].QualifiedName] = &current[i]
	}

	have := make(map[string]bool, len(sess.Functions))
	kept := sess.Functions[:0]
	for i := range sess.Functions {
		fn := &sess.Functions[i]
		c, ok := cur[fn.QualifiedName]
		if !ok {
			removed++
			continue // gone from source: drop
		}
		// Refresh positional metadata; preserve progress.
		fn.File, fn.StartLine, fn.EndLine, fn.Exported = c.File, c.StartLine, c.EndLine, c.Exported
		have[fn.QualifiedName] = true
		kept = append(kept, *fn)
	}
	sess.Functions = kept

	newNames := make(map[string]bool)
	for i := range current {
		c := &current[i]
		if have[c.QualifiedName] {
			continue // already in session, or a duplicate QualifiedName just added
		}
		nf := *c
		nf.Status = model.StatusTodo
		sess.Functions = append(sess.Functions, nf)
		have[c.QualifiedName] = true // guard duplicate QualifiedName in source
		newNames[c.QualifiedName] = true
	}

	// Re-collect new-function pointers from the final (possibly reallocated)
	// slice so callers can measure them safely.
	for i := range sess.Functions {
		if newNames[sess.Functions[i].QualifiedName] {
			added = append(added, &sess.Functions[i])
		}
	}

	sess.CheckedAt = time.Now()
	return added, removed, nil
}
