//ff:func feature=gate type=helper control=sequence lang=python
//ff:what runPyFormatter: pipes src through a stdin/stdout Python formatter (black, isort) when it is locally available, returning the formatted output; on any error (tool absent, timeout, non-zero, empty output) it returns src unchanged. The best-effort primitive tidyPySource composes — formatters are optional external tools, never required (plan §5: "black/isort(있으면). 없으면 언랩만").
package tsmagate

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

// runPyFormatter runs `name args...` with src on stdin and returns its stdout,
// or src unchanged on any failure.
func runPyFormatter(src, name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
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
