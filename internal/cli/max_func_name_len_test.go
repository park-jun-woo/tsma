package cli

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestMaxFuncNameLen_typical(t *testing.T) {
	functions := []model.Function{
		{Name: "A"},
		{Name: "LongName"},
		{Name: "Mid"},
	}
	got := maxFuncNameLen(functions)
	if got != 8 {
		t.Errorf("expected 8, got %d", got)
	}
}

func TestMaxFuncNameLen_empty(t *testing.T) {
	got := maxFuncNameLen([]model.Function{})
	if got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestMaxFuncNameLen_single(t *testing.T) {
	got := maxFuncNameLen([]model.Function{{Name: "Hello"}})
	if got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
}

func TestMaxFuncNameLen_sameLengths(t *testing.T) {
	functions := []model.Function{
		{Name: "abc"},
		{Name: "def"},
		{Name: "ghi"},
	}
	got := maxFuncNameLen(functions)
	if got != 3 {
		t.Errorf("expected 3, got %d", got)
	}
}
