package runner

import (
	"strings"
	"testing"
)

func TestUnsupportedRunnerRunReturnsError(t *testing.T) {
	r := &UnsupportedRunner{Lang: "rust"}
	result, err := r.Run("/tmp", mkMatch("main_test.rs"))
	if err == nil {
		t.Fatal("expected error from UnsupportedRunner.Run")
	}
	if result != nil {
		t.Errorf("result = %v, want nil", result)
	}
	if !strings.Contains(err.Error(), "rust") {
		t.Errorf("error should mention 'rust': %v", err)
	}
}

func TestUnsupportedRunnerRunDifferentLang(t *testing.T) {
	r := &UnsupportedRunner{Lang: "ruby"}
	_, err := r.Run("/tmp", mkMatch("test.rb"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "ruby") {
		t.Errorf("error should mention 'ruby': %v", err)
	}
}
