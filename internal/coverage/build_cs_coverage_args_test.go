package coverage

import (
	"reflect"
	"testing"
)

func TestBuildCsCoverageArgs(t *testing.T) {
	got := buildCsCoverageArgs(".tsma/coverage")
	want := []string{"test", "--collect:XPlat Code Coverage", "--results-directory", ".tsma/coverage"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
