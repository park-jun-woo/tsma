//ff:func feature=session type=implementation control=sequence
//ff:what Serializes and writes the session to disk creating directories as needed
package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Save writes the session to disk, creating directories as needed.
func Save(projectRoot string, s *model.Session) error {
	s.RecalcSummary()
	dir := Dir(projectRoot)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create session dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(dir, testDir), 0o755); err != nil {
		return fmt.Errorf("create tests dir: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}
	p := filepath.Join(dir, sessionFile)
	if err := os.WriteFile(p, data, 0o644); err != nil {
		return fmt.Errorf("write session: %w", err)
	}
	return nil
}
