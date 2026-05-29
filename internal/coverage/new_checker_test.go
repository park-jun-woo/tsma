package coverage

import "testing"

func TestNewCheckerReturnsGoChecker(t *testing.T) {
	c := NewChecker("go")
	if _, ok := c.(*GoChecker); !ok {
		t.Fatalf("NewChecker(\"go\") returned %T, want *GoChecker", c)
	}
}

func TestNewCheckerReturnsTSChecker(t *testing.T) {
	c := NewChecker("typescript")
	if _, ok := c.(*TSChecker); !ok {
		t.Fatalf("NewChecker(\"typescript\") returned %T, want *TSChecker", c)
	}
}

func TestNewCheckerReturnsPyChecker(t *testing.T) {
	c := NewChecker("python")
	if _, ok := c.(*PyChecker); !ok {
		t.Fatalf("NewChecker(\"python\") returned %T, want *PyChecker", c)
	}
}

func TestNewCheckerReturnsRsChecker(t *testing.T) {
	c := NewChecker("rust")
	if _, ok := c.(*RsChecker); !ok {
		t.Fatalf("NewChecker(\"rust\") returned %T, want *RsChecker", c)
	}
}

func TestNewCheckerReturnsJavaChecker(t *testing.T) {
	c := NewChecker("java")
	if _, ok := c.(*JavaChecker); !ok {
		t.Fatalf("NewChecker(\"java\") returned %T, want *JavaChecker", c)
	}
}

func TestNewCheckerReturnsCsChecker(t *testing.T) {
	c := NewChecker("csharp")
	if _, ok := c.(*CsChecker); !ok {
		t.Fatalf("NewChecker(\"csharp\") returned %T, want *CsChecker", c)
	}
}

func TestNewCheckerReturnsUnsupported(t *testing.T) {
	c := NewChecker("kotlin")
	u, ok := c.(*UnsupportedChecker)
	if !ok {
		t.Fatalf("NewChecker(\"kotlin\") returned %T, want *UnsupportedChecker", c)
	}
	if u.Lang != "kotlin" {
		t.Errorf("Lang = %q, want %q", u.Lang, "kotlin")
	}
}

func TestNewCheckerEmptyLang(t *testing.T) {
	c := NewChecker("")
	if _, ok := c.(*UnsupportedChecker); !ok {
		t.Fatalf("NewChecker(\"\") returned %T, want *UnsupportedChecker", c)
	}
}
