//ff:func feature=gate type=helper control=sequence lang=rust
//ff:what tidyRsSource: formats generated Rust via `rustfmt` (stdin filter) when it is locally available, otherwise returns the input unchanged. Best-effort by design (the plan §5.4: "rustfmt(있으면), 없으면 펜스 언랩만") — rustfmt is an optional external tool, never required, and any failure (absent, timeout, non-zero, empty) degrades to the unformatted-but-unwrapped source. Sibling of tidyCsSource/tidyGoSource.
package tsmagate

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

// tidyRsSource runs rustfmt over src via stdin if it is installed, returning the
// formatted output; on any error (rustfmt absent, timeout, or failure) it returns
// src unchanged.
func tidyRsSource(src string) string {
	if _, err := exec.LookPath("rustfmt"); err != nil {
		return src
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "rustfmt", "--edition", "2021", "--emit", "stdout")
	cmd.Stdin = strings.NewReader(src)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil || strings.TrimSpace(out.String()) == "" {
		return src
	}
	return out.String()
}
