package runner

import "testing"

func TestBuildGoTestArgsWithFuncs(t *testing.T) {
	args := buildGoTestArgs("./internal/handler", []string{"TestLogin", "TestSignup"})
	expected := []string{"test", "-v", "-count=1", "-run", "^(TestLogin|TestSignup)$", "./internal/handler"}

	if len(args) != len(expected) {
		t.Fatalf("got %v (len %d), want %v (len %d)", args, len(args), expected, len(expected))
	}
	for i := range args {
		if args[i] != expected[i] {
			t.Errorf("args[%d] = %q, want %q", i, args[i], expected[i])
		}
	}
}

func TestBuildGoTestArgsNoFuncs(t *testing.T) {
	args := buildGoTestArgs("./internal/handler", nil)
	expected := []string{"test", "-v", "-count=1", "./internal/handler"}

	if len(args) != len(expected) {
		t.Fatalf("got %v (len %d), want %v (len %d)", args, len(args), expected, len(expected))
	}
	for i := range args {
		if args[i] != expected[i] {
			t.Errorf("args[%d] = %q, want %q", i, args[i], expected[i])
		}
	}
}

func TestBuildGoTestArgsEmptyFuncs(t *testing.T) {
	args := buildGoTestArgs("./pkg/auth", []string{})
	expected := []string{"test", "-v", "-count=1", "./pkg/auth"}

	if len(args) != len(expected) {
		t.Fatalf("got %v (len %d), want %v (len %d)", args, len(args), expected, len(expected))
	}
	for i := range args {
		if args[i] != expected[i] {
			t.Errorf("args[%d] = %q, want %q", i, args[i], expected[i])
		}
	}
}

func TestBuildGoTestArgsSingleFunc(t *testing.T) {
	args := buildGoTestArgs("./cmd", []string{"TestMain"})
	expected := []string{"test", "-v", "-count=1", "-run", "^TestMain$", "./cmd"}

	if len(args) != len(expected) {
		t.Fatalf("got %v (len %d), want %v (len %d)", args, len(args), expected, len(expected))
	}
	for i := range args {
		if args[i] != expected[i] {
			t.Errorf("args[%d] = %q, want %q", i, args[i], expected[i])
		}
	}
}
