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

func TestNewMatcherReturnsRsMatcher(t *testing.T) {
	m := NewMatcher("rust")
	if _, ok := m.(*RsMatcher); !ok {
		t.Errorf("NewMatcher(\"rust\") returned %T, want *RsMatcher", m)
	}
}

func TestNewMatcherReturnsJavaMatcher(t *testing.T) {
	m := NewMatcher("java")
	if _, ok := m.(*JavaMatcher); !ok {
		t.Errorf("NewMatcher(\"java\") returned %T, want *JavaMatcher", m)
	}
}

func TestNewMatcherReturnsCsMatcher(t *testing.T) {
	m := NewMatcher("csharp")
	if _, ok := m.(*CsMatcher); !ok {
		t.Errorf("NewMatcher(\"csharp\") returned %T, want *CsMatcher", m)
	}
}

func TestNewMatcherReturnsUnsupported(t *testing.T) {
	m := NewMatcher("kotlin")
	_, found := m.Match("/tmp", "handler.kt")
	if found {
		t.Error("unsupported matcher should return found=false")
	}
}
