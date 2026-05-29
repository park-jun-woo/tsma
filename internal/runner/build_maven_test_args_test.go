package runner

import "testing"

func TestBuildMavenTestArgs(t *testing.T) {
	args := buildMavenTestArgs("FooTest")
	expected := []string{"-Dtest=FooTest", "test"}
	assertArgs(t, args, expected)
}
