//ff:func feature=runner type=helper control=sequence
//ff:what Checks project configuration files to determine if pytest is used
package runner

import "path/filepath"

// detectPytest checks if the project uses pytest.
func detectPytest(projectRoot string) bool {
	if containsPytest(filepath.Join(projectRoot, "pyproject.toml"), "[tool.pytest") {
		return true
	}
	if containsPytest(filepath.Join(projectRoot, "setup.cfg"), "[tool:pytest]") {
		return true
	}
	if containsPytest(filepath.Join(projectRoot, "requirements.txt"), "pytest") {
		return true
	}
	if containsPytest(filepath.Join(projectRoot, "requirements-dev.txt"), "pytest") {
		return true
	}
	if fileExists(filepath.Join(projectRoot, "pytest.ini")) {
		return true
	}
	if fileExists(filepath.Join(projectRoot, "conftest.py")) {
		return true
	}
	return false
}
