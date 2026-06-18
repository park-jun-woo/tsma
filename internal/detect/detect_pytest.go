//ff:func feature=detect type=helper control=sequence lang=python
//ff:what Single source of truth: reports whether a project uses pytest
package detect

import "path/filepath"

// DetectPytest is the single source of truth for whether a Python project uses
// pytest. Both the runner and the coverage stages must agree through this one
// function (cf. BUG-001 / Phase006: a runner that fell back to unittest while
// coverage already assumed pytest).
//
// It checks, in cheapest-first order:
//   - legacy/static signals: [tool.pytest] in pyproject.toml, [tool:pytest] in
//     setup.cfg, a pytest token in requirements*.txt, or pytest.ini/conftest.py;
//   - PEP 621 dependency declarations: pytest under
//     [project.optional-dependencies] / [dependency-groups] / [project]
//     dependencies (containsPytestDep);
//   - conventional layout: tests/ holding test_*.py / *_test.py (hasPytestLayout);
//   - last resort: a project-local venv carrying bin/pytest (probePytestVenv).
//
// Any single signal is sufficient. A pure unittest project (no pytest
// declaration) trips none of these and stays false.
func DetectPytest(projectRoot string) bool {
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

	// D1 — extended signals (PEP 621 declarations, layout, venv probe).
	if containsPytestDep(projectRoot) {
		return true
	}
	if hasPytestLayout(projectRoot) {
		return true
	}
	if probePytestVenv(projectRoot) {
		return true
	}
	return false
}
