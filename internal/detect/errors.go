package detect

import "errors"

// ErrUnsupportedLanguage is returned when no supported language is detected.
var ErrUnsupportedLanguage = errors.New("unsupported language: no go.mod, package.json, or pyproject.toml found")
