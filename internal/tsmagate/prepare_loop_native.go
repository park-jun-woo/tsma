//ff:func feature=gate type=helper control=selection
//ff:what prepareLoopNative: dispatches loop-mode measurement to the language's native non-invasive pipeline — Go via overlay (prepareLoopGo), TypeScript via .tsma/test isolation (prepareLoopTS). Returns (nil, false) for languages without a native loop path so Prepare falls through to the generic disk-write loop. Keeps Prepare's two early-return branches collapsed into one (length + nesting budget).
package tsmagate

import "github.com/park-jun-woo/reins/pkg/quest"

// prepareLoopNative runs the per-language non-invasive loop measurement and
// returns (measurement, true) when the language has one; otherwise (nil, false).
func prepareLoopNative(it *quest.Item, p funcPayload, raw []byte) (*measurement, bool) {
	switch p.Lang {
	case "go":
		return prepareLoopGo(it, p, raw), true
	case "typescript":
		return prepareLoopTS(it, p, raw), true
	}
	return nil, false
}
