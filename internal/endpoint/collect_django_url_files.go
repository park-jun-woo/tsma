//ff:func feature=endpoint type=implementation control=sequence
//ff:what Walks project to find Django urls.py files
package endpoint

import (
	"fmt"
	"os"
	"path/filepath"
)

func collectDjangoURLFiles(projectRoot string) ([]string, error) {
	var urlsFiles []string

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if djangoSkipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Base(path) == "urls.py" {
			urlsFiles = append(urlsFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk project: %w", err)
	}

	return urlsFiles, nil
}
