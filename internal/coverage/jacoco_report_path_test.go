package coverage

import (
	"path/filepath"
	"testing"
)

func TestJacocoReportPathMaven(t *testing.T) {
	got := jacocoReportPath("/proj", "maven")
	want := filepath.Join("/proj", "target", "site", "jacoco", "jacoco.xml")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestJacocoReportPathGradle(t *testing.T) {
	got := jacocoReportPath("/proj", "gradle")
	want := filepath.Join("/proj", "build", "reports", "jacoco", "test", "jacocoTestReport.xml")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestJacocoReportPathUnknown(t *testing.T) {
	if got := jacocoReportPath("/proj", "ant"); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}
