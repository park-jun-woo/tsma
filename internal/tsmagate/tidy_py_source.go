//ff:func feature=gate type=helper control=sequence lang=python
//ff:what tidyPySource: formats a generated Python test best-effort — isort to order imports, then black to normalize style — each via runPyFormatter (stdin/stdout), each independently skipped when its tool is absent. The Python analogue of tidyTSSource; never required, so a clean environment degrades to the unformatted-but-unwrapped source.
package tsmagate

// tidyPySource runs isort then black on src when available, returning the
// formatted output; missing tools are silently skipped.
func tidyPySource(src string) string {
	src = runPyFormatter(src, "isort", "-")
	src = runPyFormatter(src, "black", "-q", "-")
	return src
}
