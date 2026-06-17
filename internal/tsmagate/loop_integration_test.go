//ff:func feature=gate type=test
//ff:what loop 통합테스트(§5 시나리오 1~4): 실제 cli.NewQuestCmd 트리에 stub backend(llm.CallFunc)를 주입해 scan→loop를 구동한다. 1) 완전커버 1회 PASS, 2) 부분→완전 2회 수렴, 3) 컴파일깨짐 MaxTries DONE(무한루프 없음), 4) --max-items 1로 1개만 처리. 네트워크/claude 서브프로세스 없이 결정론적(go test 서브프로세스만, 기존 prepare_test와 동일).

package tsmagate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/cli"
	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// runTsma builds a fresh tsma command tree (NewQuestCmd) with the given options
// and runs one invocation against the session/out paths, returning combined
// stdout+stderr. A fresh tree per call mirrors separate process invocations
// (scan then loop) while sharing the on-disk session.
func runTsma(t *testing.T, opts cli.Options, session, out string, args ...string) (string, error) {
	t.Helper()
	root := cli.NewQuestCmd("tsma", New(), opts)
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	full := append(append([]string{}, args...), "--session", session, "--out", out)
	root.SetArgs(full)
	err := root.Execute()
	return buf.String(), err
}

// loopOpts returns the production LoopOptions with the stub backend injected, so
// the test exercises the real C2 config (System/RuleSystem) but never touches a
// network or the claude CLI.
func loopOpts(backend llm.Backend) cli.Options {
	lo := LoopOptions()
	lo.LLM = backend
	return cli.Options{Loop: lo}
}

const classifySrc = "package pkg\n\nfunc Classify(n int) string {\n\tif n > 0 {\n\t\treturn \"pos\"\n\t}\n\treturn \"nonpos\"\n}\n"

const classifyFullTest = "package pkg\n\nimport \"testing\"\n\n" +
	"func TestClassify(t *testing.T) {\n" +
	"\tif Classify(1) != \"pos\" {\n\t\tt.Fatal(\"pos\")\n\t}\n" +
	"\tif Classify(-1) != \"nonpos\" {\n\t\tt.Fatal(\"nonpos\")\n\t}\n}\n"

const classifyPartialTest = "package pkg\n\nimport \"testing\"\n\n" +
	"func TestClassify(t *testing.T) {\n" +
	"\tif Classify(1) != \"pos\" {\n\t\tt.Fatal(\"pos\")\n\t}\n}\n"

// classifyModule writes a minimal Go module with one branching function.
func classifyModule(t *testing.T) string {
	return writeGoPkg(t, map[string]string{
		"go.mod":          "module classifymod\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
}

// sessionPaths returns separate session/out paths in their own temp dir so they
// never land in the fixture module (keeping scan's indexing input clean).
func sessionPaths(t *testing.T) (string, string) {
	d := t.TempDir()
	return d + "/session.json", d + "/out.jsonl"
}

// Scenario 1: a complete covering test in one shot → PASS locked on first try.
func TestLoop_FullCoverageOneShot(t *testing.T) {
	root := classifyModule(t)
	chdirTo(t, root) // runner/coverage resolve relative test paths against cwd
	session, out := sessionPaths(t)
	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		return classifyFullTest, nil
	})
	opts := loopOpts(backend)

	if _, err := runTsma(t, opts, session, out, "scan", root); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := runTsma(t, opts, session, out, "loop")
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if calls != 1 {
		t.Fatalf("backend called %d times, want 1 (one-shot PASS)", calls)
	}
	if !strings.Contains(got, "-> PASS") {
		t.Fatalf("loop output missing PASS:\n%s", got)
	}
}

// Scenario 2: a partial test (one branch) then a full test → FAIL then PASS,
// proving the feedback loop converges (coverage drives the retry).
func TestLoop_PartialThenFullConverges(t *testing.T) {
	root := classifyModule(t)
	chdirTo(t, root)
	session, out := sessionPaths(t)
	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		if calls == 1 {
			return classifyPartialTest, nil // passes but coverage < 100
		}
		return classifyFullTest, nil // retry covers the missing branch
	})
	opts := loopOpts(backend)

	if _, err := runTsma(t, opts, session, out, "scan", root); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := runTsma(t, opts, session, out, "loop")
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if calls != 2 {
		t.Fatalf("backend called %d times, want 2 (FAIL then PASS)", calls)
	}
	if !strings.Contains(got, "-> FAIL") || !strings.Contains(got, "-> PASS") {
		t.Fatalf("loop output = %q, want a FAIL then a PASS", got)
	}
}

// Scenario 3: a compile-broken test every attempt → tests-must-pass FAIL each
// try until MaxTries locks the item DONE. The loop terminates (no infinite spin).
func TestLoop_BrokenCompileLocksDoneAtMaxTries(t *testing.T) {
	root := classifyModule(t)
	chdirTo(t, root)
	session, out := sessionPaths(t)
	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		return "package pkg\n\nthis is not valid go\n", nil
	})
	opts := loopOpts(backend)

	if _, err := runTsma(t, opts, session, out, "scan", root); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := runTsma(t, opts, session, out, "loop")
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if calls != quest.MaxTries {
		t.Fatalf("backend called %d times, want %d (one per try until DONE)", calls, quest.MaxTries)
	}
	if strings.Contains(got, "-> PASS") {
		t.Fatalf("a broken test must never PASS:\n%s", got)
	}
	if !strings.Contains(got, "state DONE") {
		t.Fatalf("expected the item to lock DONE at MaxTries:\n%s", got)
	}
}

// Scenario 4: --max-items 1 processes exactly one TODO even when more remain.
func TestLoop_MaxItemsCapsWork(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":       "module capmod\n\ngo 1.22\n",
		"pkg/alpha.go": "package pkg\n\nfunc Alpha() int { return 1 }\n",
		"pkg/beta.go":  "package pkg\n\nfunc Beta() int { return 2 }\n",
	})
	chdirTo(t, root)
	session, out := sessionPaths(t)
	const alphaTest = "package pkg\n\nimport \"testing\"\n\nfunc TestAlpha(t *testing.T) { if Alpha() != 1 { t.Fatal(\"x\") } }\n"
	const betaTest = "package pkg\n\nimport \"testing\"\n\nfunc TestBeta(t *testing.T) { if Beta() != 2 { t.Fatal(\"x\") } }\n"
	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		if strings.Contains(user, "Alpha") {
			return alphaTest, nil
		}
		return betaTest, nil
	})
	opts := loopOpts(backend)

	if _, err := runTsma(t, opts, session, out, "scan", root); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := runTsma(t, opts, session, out, "loop", "--max-items", "1")
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if calls != 1 {
		t.Fatalf("backend called %d times, want 1 (--max-items 1)", calls)
	}
	if !strings.Contains(got, "processed 1 item(s)") {
		t.Fatalf("expected 'processed 1 item(s)':\n%s", got)
	}
}
