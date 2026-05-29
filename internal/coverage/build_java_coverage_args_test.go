package coverage

import (
	"reflect"
	"testing"
)

func TestBuildJavaCoverageArgsMaven(t *testing.T) {
	got := buildJavaCoverageArgs("maven")
	want := []string{"test", "jacoco:report"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildJavaCoverageArgsGradle(t *testing.T) {
	got := buildJavaCoverageArgs("gradle")
	want := []string{"test", "jacocoTestReport"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildJavaCoverageArgsUnknown(t *testing.T) {
	if got := buildJavaCoverageArgs("ant"); got != nil {
		t.Errorf("got %v, want nil", got)
	}
}
