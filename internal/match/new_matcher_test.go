package match

import "testing"

func TestNewMatcherReturnsGoMatcher(t *testing.T) {
	m := NewMatcher("go")
	if _, ok := m.(*GoMatcher); !ok {
		t.Errorf("NewMatcher(\"go\") returned %T, want *GoMatcher", m)
	}
}

func TestNewMatcherReturnsTSMatcher(t *testing.T) {
	m := NewMatcher("typescript")
	if _, ok := m.(*TSMatcher); !ok {
		t.Errorf("NewMatcher(\"typescript\") returned %T, want *TSMatcher", m)
	}
}

func TestNewMatcherReturnsPyMatcher(t *testing.T) {
	m := NewMatcher("python")
	if _, ok := m.(*PyMatcher); !ok {
		t.Errorf("NewMatcher(\"python\") returned %T, want *PyMatcher", m)
	}
}

func TestNewMatcherReturnsUnsupported(t *testing.T) {
	m := NewMatcher("rust")
	_, found := m.Match("/tmp", "handler.rs")
	if found {
		t.Error("unsupported matcher should return found=false")
	}
}
