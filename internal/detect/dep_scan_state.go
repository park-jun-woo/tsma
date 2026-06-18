//ff:type feature=detect type=model lang=python
//ff:what Carries the line-scan state for PEP 621 dependency detection
package detect

// depScanState tracks where a pyproject.toml line scan currently is while
// looking for a pytest dependency declaration (see advanceDepScan).
type depScanState struct {
	inDepTable bool // inside [project.optional-dependencies] / [dependency-groups]
	inDepArray bool // inside a multiline `dependencies = [ ... ]` array
}
