package detect

import (
	"errors"
	"testing"
)

func TestErrUnsupportedLanguageNotNil(t *testing.T) {
	if ErrUnsupportedLanguage == nil {
		t.Fatal("ErrUnsupportedLanguage is nil")
	}
}

func TestErrUnsupportedLanguageMessage(t *testing.T) {
	msg := ErrUnsupportedLanguage.Error()
	if msg == "" {
		t.Error("ErrUnsupportedLanguage.Error() is empty")
	}
}

func TestErrUnsupportedLanguageIsComparable(t *testing.T) {
	err := ErrUnsupportedLanguage
	if !errors.Is(err, ErrUnsupportedLanguage) {
		t.Error("errors.Is should match ErrUnsupportedLanguage")
	}
}
