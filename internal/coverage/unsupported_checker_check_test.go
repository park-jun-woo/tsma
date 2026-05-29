package coverage

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestUnsupportedCheckerCheck(t *testing.T) {
	checker := &UnsupportedChecker{Lang: "ruby"}
	fn := &model.Function{
		Name:      "TestFunc",
		File:      "test.rb",
		StartLine: 1,
		EndLine:   10,
	}

	report, err := checker.Check("/project", mkMatch("test_test.rb"), fn)
	if report != nil {
		t.Error("expected nil report from UnsupportedChecker")
	}
	if err == nil {
		t.Fatal("expected error from UnsupportedChecker")
	}

	unsupErr, ok := err.(*ErrUnsupported)
	if !ok {
		t.Fatalf("error type = %T, want *ErrUnsupported", err)
	}
	if unsupErr.Lang != "ruby" {
		t.Errorf("Lang = %q, want %q", unsupErr.Lang, "ruby")
	}
}

func TestUnsupportedCheckerCheckEmptyLang(t *testing.T) {
	checker := &UnsupportedChecker{Lang: ""}
	_, err := checker.Check("", mkMatch(""), &model.Function{})
	if err == nil {
		t.Fatal("expected error from UnsupportedChecker even with empty lang")
	}
}
