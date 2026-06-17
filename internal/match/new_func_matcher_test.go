package match

import "testing"

// TestNewFuncMatcherGo verifies the Go language returns the content-aware
// GoFuncMatcher.
func TestNewFuncMatcherGo(t *testing.T) {
	if _, ok := NewFuncMatcher("go").(*GoFuncMatcher); !ok {
		t.Fatalf("NewFuncMatcher(go) = %T, want *GoFuncMatcher", NewFuncMatcher("go"))
	}
}

// TestNewFuncMatcherTypescript verifies TypeScript returns the content-aware
// TypeScriptFuncMatcher (Phase005a D2).
func TestNewFuncMatcherTypescript(t *testing.T) {
	if _, ok := NewFuncMatcher("typescript").(*TypeScriptFuncMatcher); !ok {
		t.Fatalf("NewFuncMatcher(typescript) = %T, want *TypeScriptFuncMatcher", NewFuncMatcher("typescript"))
	}
}

// TestNewFuncMatcherPython verifies Python returns the content-aware
// PythonFuncMatcher (Phase005b D2).
func TestNewFuncMatcherPython(t *testing.T) {
	if _, ok := NewFuncMatcher("python").(*PythonFuncMatcher); !ok {
		t.Fatalf("NewFuncMatcher(python) = %T, want *PythonFuncMatcher", NewFuncMatcher("python"))
	}
}

// TestNewFuncMatcherJava verifies Java returns the content-aware JavaFuncMatcher
// (Phase005c D2; JavaMatcher is retained as its last-resort filename fallback).
func TestNewFuncMatcherJava(t *testing.T) {
	if _, ok := NewFuncMatcher("java").(*JavaFuncMatcher); !ok {
		t.Fatalf("NewFuncMatcher(java) = %T, want *JavaFuncMatcher", NewFuncMatcher("java"))
	}
}

// TestNewFuncMatcherNonGoFallback verifies that every remaining non-content
// language returns a fallback adapter wrapping the language's legacy Matcher.
func TestNewFuncMatcherNonGoFallback(t *testing.T) {
	cases := map[string]Matcher{
		"rust":    &RsMatcher{},
		"csharp":  &CsMatcher{},
		"unknown": &unsupportedMatcher{},
	}
	for lang, wantInner := range cases {
		fm := NewFuncMatcher(lang)
		adapter, ok := fm.(*fallbackFuncMatcher)
		if !ok {
			t.Errorf("NewFuncMatcher(%q) = %T, want *fallbackFuncMatcher", lang, fm)
			continue
		}
		gotType := typeName(adapter.inner)
		wantType := typeName(wantInner)
		if gotType != wantType {
			t.Errorf("NewFuncMatcher(%q) inner = %s, want %s", lang, gotType, wantType)
		}
	}
}

func typeName(m Matcher) string {
	switch m.(type) {
	case *TSMatcher:
		return "TSMatcher"
	case *PyMatcher:
		return "PyMatcher"
	case *RsMatcher:
		return "RsMatcher"
	case *JavaMatcher:
		return "JavaMatcher"
	case *CsMatcher:
		return "CsMatcher"
	case *unsupportedMatcher:
		return "unsupportedMatcher"
	default:
		return "unknown"
	}
}
