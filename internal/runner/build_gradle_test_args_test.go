package runner

import "testing"

func TestBuildGradleTestArgs(t *testing.T) {
	args := buildGradleTestArgs("FooTest")
	expected := []string{"test", "--tests", "FooTest"}
	assertArgs(t, args, expected)
}
