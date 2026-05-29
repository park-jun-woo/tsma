package index

import "testing"

func TestMatchPatternTrailingSlashDir(t *testing.T) {
	// Pattern "vendor/" should match a directory named "vendor"
	if !matchPattern("vendor", "vendor", true, "vendor/") {
		t.Error("expected vendor/ to match directory vendor")
	}
	// Should not match a file
	if matchPattern("vendor", "vendor", false, "vendor/") {
		t.Error("expected vendor/ to not match file vendor")
	}
}

func TestMatchPatternTrailingSlashNestedDir(t *testing.T) {
	// Pattern "build/" should match nested directory
	if !matchPattern("src/build", "build", true, "build/") {
		t.Error("expected build/ to match nested directory src/build")
	}
}

func TestMatchPatternTrailingSlashMidPathSegment(t *testing.T) {
	// dirPattern appears as an interior path segment (not a suffix and the
	// directory name itself differs) -> exercises the Contains() branch.
	if !matchPattern("a/build/c", "c", true, "build/") {
		t.Error("expected build/ to match interior segment of a/build/c")
	}
}

func TestMatchPatternTrailingSlashNoMatch(t *testing.T) {
	// A directory whose name/path neither equals, suffixes, nor contains the
	// dir pattern segment -> all dir branches return false.
	if matchPattern("src/lib", "lib", true, "build/") {
		t.Error("expected build/ to not match src/lib")
	}
}

func TestMatchPatternWithSlash(t *testing.T) {
	// Pattern "internal/tmp/*.go" should match path glob
	if !matchPattern("internal/tmp/foo.go", "foo.go", false, "internal/tmp/*.go") {
		t.Error("expected internal/tmp/*.go to match internal/tmp/foo.go")
	}
	if matchPattern("internal/other/foo.go", "foo.go", false, "internal/tmp/*.go") {
		t.Error("expected internal/tmp/*.go to not match internal/other/foo.go")
	}
}

func TestMatchPatternWithSlashDotSlashPrefix(t *testing.T) {
	// Paths starting with "./" should be trimmed before matching
	if !matchPattern("./internal/tmp/foo.go", "foo.go", false, "internal/tmp/*.go") {
		t.Error("expected internal/tmp/*.go to match ./internal/tmp/foo.go")
	}
}

func TestMatchPatternFileGlob(t *testing.T) {
	// Pattern without slash should match filename only
	if !matchPattern("src/app.log", "app.log", false, "*.log") {
		t.Error("expected *.log to match app.log")
	}
	if matchPattern("src/app.txt", "app.txt", false, "*.log") {
		t.Error("expected *.log to not match app.txt")
	}
}

func TestMatchPatternExactFilename(t *testing.T) {
	if !matchPattern("Makefile", "Makefile", false, "Makefile") {
		t.Error("expected Makefile pattern to match Makefile")
	}
	if matchPattern("README.md", "README.md", false, "Makefile") {
		t.Error("expected Makefile pattern to not match README.md")
	}
}
