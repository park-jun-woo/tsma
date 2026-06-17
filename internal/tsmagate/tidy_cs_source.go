//ff:func feature=gate type=helper control=sequence lang=csharp
//ff:what tidyCsSource: formats generated C# via `dotnet format` when the .NET SDK is locally available, otherwise returns the input unchanged. Because dotnet format is file/folder-scoped (not a stdin filter like gofmt), the source is written to a temp .cs file, formatted in place with `dotnet format whitespace --folder`, and read back. Best-effort by design (the plan: "dotnet format(있으면), 없으면 언랩만") — the SDK is an optional external tool, never required, and any failure (absent, timeout, non-zero, empty) degrades to the unformatted-but-unwrapped source.
package tsmagate

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// tidyCsSource runs `dotnet format` on src if the SDK is installed, returning the
// formatted output; on any error (SDK absent, timeout, or failure) it returns src
// unchanged.
func tidyCsSource(src string) string {
	if _, err := exec.LookPath("dotnet"); err != nil {
		return src
	}

	dir, err := os.MkdirTemp("", "tsma-cs-fmt-")
	if err != nil {
		return src
	}
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "Generated.cs")
	if err := os.WriteFile(file, []byte(src), 0o644); err != nil {
		return src
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "dotnet", "format", "whitespace", "--folder", dir)
	if err := cmd.Run(); err != nil {
		return src
	}

	formatted, err := os.ReadFile(file)
	if err != nil || strings.TrimSpace(string(formatted)) == "" {
		return src
	}
	return string(formatted)
}
