//ff:func feature=gate type=test
//ff:what test helpers for driving the best-effort external-tool branches without
// the real tools: installFakeTool writes an executable shell script into a fresh
// temp dir and prepends it to PATH (t.Setenv) so exec.LookPath/exec.Command in
// the code under test resolve the fake — making the dotnet/rustfmt/
// google-java-format/npx success, failure, and empty-output branches
// deterministically reachable in any environment.
package tsmagate

import (
	"os"
	"path/filepath"
	"testing"
)

// installFakeTool writes an executable shell script named name into a fresh
// temp dir and prepends that dir to PATH, so the code under test resolves the
// fake instead of any real tool. Returns the script's absolute path.
func installFakeTool(t *testing.T, name, script string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return path
}

// emptyPath points PATH at a single empty temp dir so LookPath deterministically
// fails for every tool (the formatter-absent branch), even on machines where the
// real tool is installed.
func emptyPath(t *testing.T) {
	t.Helper()
	t.Setenv("PATH", t.TempDir())
}
