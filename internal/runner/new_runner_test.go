package runner

import "testing"

func TestNewRunnerReturnsGoRunner(t *testing.T) {
	r := NewRunner("go")
	if _, ok := r.(*GoRunner); !ok {
		t.Errorf("NewRunner(\"go\") returned %T, want *GoRunner", r)
	}
}

func TestNewRunnerReturnsTSRunner(t *testing.T) {
	r := NewRunner("typescript")
	if _, ok := r.(*TSRunner); !ok {
		t.Errorf("NewRunner(\"typescript\") returned %T, want *TSRunner", r)
	}
}

func TestNewRunnerReturnsPyRunner(t *testing.T) {
	r := NewRunner("python")
	if _, ok := r.(*PyRunner); !ok {
		t.Errorf("NewRunner(\"python\") returned %T, want *PyRunner", r)
	}
}

func TestNewRunnerReturnsUnsupported(t *testing.T) {
	r := NewRunner("rust")
	u, ok := r.(*UnsupportedRunner)
	if !ok {
		t.Fatalf("NewRunner(\"rust\") returned %T, want *UnsupportedRunner", r)
	}
	if u.Lang != "rust" {
		t.Errorf("Lang = %q, want %q", u.Lang, "rust")
	}
}

func TestNewRunnerEmptyString(t *testing.T) {
	r := NewRunner("")
	if _, ok := r.(*UnsupportedRunner); !ok {
		t.Errorf("NewRunner(\"\") returned %T, want *UnsupportedRunner", r)
	}
}

func TestNewRunnerImplementsInterface(t *testing.T) {
	langs := []string{"go", "typescript", "python", "ruby"}
	for _, lang := range langs {
		r := NewRunner(lang)
		// All returned values should satisfy the Runner interface
		var _ Runner = r
	}
}
