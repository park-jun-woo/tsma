//ff:func feature=gate type=helper control=sequence lang=typescript
//ff:what tidyTSSource: formats generated TS via prettier when it is locally available (`npx --no-install prettier --stdin-filepath`), otherwise returns the input unchanged. Best-effort by design (the plan: "prettier(있으면). 없으면 언랩만") — prettier is an optional external tool, never required, and any failure (absent, timeout, non-zero) degrades to the unformatted-but-unwrapped source.
package tsmagate

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

// tidyTSSource runs prettier on src if it is installed locally, returning the
// formatted output; on any error (prettier absent, timeout, or failure) it
// returns src unchanged. The --no-install flag prevents npx from attempting a
// network download.
func tidyTSSource(src string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "npx", "--no-install", "prettier", "--stdin-filepath", "tsma_gen.test.ts")
	cmd.Stdin = strings.NewReader(src)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return src
	}
	formatted := out.String()
	if strings.TrimSpace(formatted) == "" {
		return src
	}
	return formatted
}
