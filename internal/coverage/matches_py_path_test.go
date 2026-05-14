package coverage

import "testing"

func TestMatchesPyPathDirectMatch(t *testing.T) {
	if !matchesPyPath("handler.py", "handler.py", "/project") {
		t.Error("expected true for direct match")
	}
}

func TestMatchesPyPathAbsoluteMatch(t *testing.T) {
	if !matchesPyPath("/project/handler.py", "handler.py", "/project") {
		t.Error("expected true for absolute path match")
	}
}

func TestMatchesPyPathSuffixWithSep(t *testing.T) {
	if !matchesPyPath("/deep/path/handler.py", "handler.py", "/project") {
		t.Error("expected true for suffix with separator")
	}
}

func TestMatchesPyPathNoMatchDifferentFile(t *testing.T) {
	if matchesPyPath("other.py", "handler.py", "/project") {
		t.Error("expected false for different file")
	}
}

func TestMatchesPyPathFalsePositive(t *testing.T) {
	if matchesPyPath("other_handler.py", "handler.py", "/project") {
		t.Error("expected false: other_handler.py should not match handler.py")
	}
}

func TestMatchesPyPathSubdirAbsMatch(t *testing.T) {
	if !matchesPyPath("/project/src/handler.py", "src/handler.py", "/project") {
		t.Error("expected true for subdirectory absolute match")
	}
}

func TestMatchesPyPathDifferentDirs(t *testing.T) {
	if matchesPyPath("other/handler.py", "src/handler.py", "/project") {
		t.Error("expected false for different directories")
	}
}
