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

// TestNewFuncMatcherNonGoFallback verifies that every remaining non-content
// language returns a fallback adapter wrapping the language's legacy Matcher.
func TestNewFuncMatcherNonGoFallback(t *testing.T) {
	cases := map[string]Matcher{
		"python": &PyMatcher{},
		"rust":       &RsMatcher{},
		"java":       &JavaMatcher{},
		"csharp":     &CsMatcher{},
		"unknown":    &unsupportedMatcher{},
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
