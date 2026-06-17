//ff:func feature=gate type=command control=sequence level=error
//ff:what Seed: 입력 디렉터리에서 언어를 판별(detect.Detect)하고 함수 목록을 인덱싱(index.NewIndexer.Index)해 함수마다 quest.Item을 만든다. Item.Key=QualifiedName, payload에 lang/root/함수위치를 스냅샷(Status/Attempt 제외 — reins Item.State/Tries 소유). rulebook E-* 인덱싱 규칙은 인덱서 내부에 그대로 보존된다.

package tsmagate

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/detect"
	"github.com/park-jun-woo/tsma/internal/index"
)

// Seed turns a project directory into the initial TODO items: one quest.Item per
// indexable function. It detects the project language, runs the language indexer
// (which applies the rulebook E-* eligibility filters internally), and snapshots
// each function into an item payload. args[0] is the project root (default ".").
// The root is resolved to an absolute path so later subcommands re-measure the
// same tree regardless of the working directory.
func (d *Definition) Seed(args []string) ([]*quest.Item, error) {
	root := "."
	if len(args) > 0 && args[0] != "" {
		root = args[0]
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve root %q: %w", root, err)
	}
	root = abs

	lf, err := detect.Detect(root)
	if err != nil {
		return nil, fmt.Errorf("detect language in %q: %w", root, err)
	}
	lang := lf.Lang

	fns, err := index.NewIndexer(lang).Index(root)
	if err != nil {
		return nil, fmt.Errorf("index %s functions in %q: %w", lang, root, err)
	}

	items := make([]*quest.Item, 0, len(fns))
	for i := range fns {
		fn := fns[i]
		// Strip the runtime fields reins now owns so the payload is pure
		// location data; State/Tries are the ratchet's, not the payload's.
		fn.Status = ""
		fn.Attempt = 0
		fn.CoveragePct = 0
		fn.TestMtime = ""
		fn.FailOutput = ""

		it := &quest.Item{Key: fn.QualifiedName, State: quest.TODO}
		if err := it.SetPayload(funcPayload{Lang: lang, Root: root, Fn: fn}); err != nil {
			return nil, fmt.Errorf("snapshot payload for %s: %w", fn.QualifiedName, err)
		}
		items = append(items, it)
	}
	return items, nil
}
