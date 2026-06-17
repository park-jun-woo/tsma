//ff:func feature=gate type=helper control=sequence level=error lang=rust
//ff:what restoreRsSource: rolls the source file back to its pre-injection contents by copying the backing (orig-<item>.rs under .tsma/test) over it — the rollback half of the D5 in-file strategy (③). Invoked by finalizeRsBacking on every non-terminal result and (via defer) on a panic, so the measurement-time mutation of the real source never leaks past the measurement. A read failure (missing backing) leaves the source as-is and reports the error to the caller, which records it on the measurement (never silent).
package tsmagate

import (
	"os"
	"path/filepath"
)

// restoreRsSource copies the backing back over the source file. It returns an
// error if the backing cannot be read (the source is then left untouched).
func restoreRsSource(p funcPayload, backingRel string) error {
	data, err := os.ReadFile(filepath.Join(p.Root, backingRel))
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(p.Root, p.Fn.File), data, 0o644)
}
