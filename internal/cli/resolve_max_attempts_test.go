package cli

import (
	"strconv"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/spf13/cobra"
)

// mkMaxAttemptsCmd builds a cobra command with the --max-attempts flag,
// optionally marking it as explicitly set with the given value.
func mkMaxAttemptsCmd(t *testing.T, set bool, val int) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{Use: "next"}
	cmd.Flags().Int("max-attempts", defaultMaxAttempts, "")
	if set {
		if err := cmd.Flags().Set("max-attempts", strconv.Itoa(val)); err != nil {
			t.Fatalf("set flag: %v", err)
		}
	}
	return cmd
}

func TestResolveMaxAttempts_flagWins(t *testing.T) {
	sess := &model.Session{MaxAttempts: 7}
	cmd := mkMaxAttemptsCmd(t, true, 5)
	if err := resolveMaxAttempts(cmd, sess); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if sess.MaxAttempts != 5 {
		t.Fatalf("flag should win: got %d, want 5", sess.MaxAttempts)
	}
}

func TestResolveMaxAttempts_sessionWhenNoFlag(t *testing.T) {
	sess := &model.Session{MaxAttempts: 7}
	cmd := mkMaxAttemptsCmd(t, false, 0)
	if err := resolveMaxAttempts(cmd, sess); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if sess.MaxAttempts != 7 {
		t.Fatalf("stored session value should be kept: got %d, want 7", sess.MaxAttempts)
	}
}

func TestResolveMaxAttempts_defaultWhenUnset(t *testing.T) {
	sess := &model.Session{MaxAttempts: 0}
	cmd := mkMaxAttemptsCmd(t, false, 0)
	if err := resolveMaxAttempts(cmd, sess); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if sess.MaxAttempts != defaultMaxAttempts {
		t.Fatalf("should fall back to default: got %d, want %d", sess.MaxAttempts, defaultMaxAttempts)
	}
}

func TestResolveMaxAttempts_rejectBelowOne(t *testing.T) {
	for _, n := range []int{0, -1, -5} {
		sess := &model.Session{}
		cmd := mkMaxAttemptsCmd(t, true, n)
		if err := resolveMaxAttempts(cmd, sess); err == nil {
			t.Fatalf("expected error for --max-attempts %d, got nil", n)
		}
	}
}
