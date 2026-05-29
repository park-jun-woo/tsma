package runner

import "testing"

func TestBuildCsTestArgs(t *testing.T) {
	args := buildCsTestArgs("FooTests")
	assertArgs(t, args, []string{"test", "--filter", "FullyQualifiedName~FooTests"})
}
