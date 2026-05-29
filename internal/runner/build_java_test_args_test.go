package runner

import "testing"

func TestBuildJavaTestArgsMaven(t *testing.T) {
	args := buildJavaTestArgs("maven", "FooTest")
	assertArgs(t, args, []string{"-Dtest=FooTest", "test"})
}

func TestBuildJavaTestArgsGradle(t *testing.T) {
	args := buildJavaTestArgs("gradle", "FooTest")
	assertArgs(t, args, []string{"test", "--tests", "FooTest"})
}

func TestBuildJavaTestArgsUnknown(t *testing.T) {
	if args := buildJavaTestArgs("ant", "FooTest"); args != nil {
		t.Errorf("expected nil for unknown build tool, got %v", args)
	}
}
