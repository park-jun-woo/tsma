//ff:func feature=gate type=helper control=sequence lang=java
//ff:what tidyJavaSource: formats generated Java via google-java-format when it is locally available (`google-java-format -` reading stdin), otherwise returns the input unchanged. Best-effort by design (the plan: "google-java-format(있으면), 없으면 언랩만") — the formatter is an optional external tool, never required, and any failure (absent, timeout, non-zero) degrades to the unformatted-but-unwrapped source.
package tsmagate

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

// tidyJavaSource runs google-java-format on src if it is installed, returning
// the formatted output; on any error (formatter absent, timeout, or failure) it
// returns src unchanged.
func tidyJavaSource(src string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "google-java-format", "-")
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
