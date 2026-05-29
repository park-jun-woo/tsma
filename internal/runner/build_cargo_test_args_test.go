package runner

import "testing"

func TestBuildCargoTestArgsIntegration(t *testing.T) {
	args := buildCargoTestArgs("tests/api.rs")
	expected := []string{"test", "--test", "api"}
	assertArgs(t, args, expected)
}

func TestBuildCargoTestArgsInFile(t *testing.T) {
	args := buildCargoTestArgs("src/lib.rs")
	expected := []string{"test"}
	assertArgs(t, args, expected)
}

func TestBuildCargoTestArgsNestedSource(t *testing.T) {
	// Unit test inside a nested source module: plain cargo test.
	args := buildCargoTestArgs("src/util/math.rs")
	expected := []string{"test"}
	assertArgs(t, args, expected)
}

func TestBuildCargoTestArgsTestsRoot(t *testing.T) {
	args := buildCargoTestArgs("tests/integration.rs")
	expected := []string{"test", "--test", "integration"}
	assertArgs(t, args, expected)
}
