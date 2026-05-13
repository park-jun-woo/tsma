//ff:func feature=session type=implementation control=sequence
//ff:what Reads and deserializes the session from disk
package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Load reads the session from disk.
func Load(projectRoot string) (*model.Session, error) {
	p := filepath.Join(Dir(projectRoot), sessionFile)
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("read session: %w", err)
	}
	var s model.Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse session: %w", err)
	}
	return &s, nil
}
