//ff:type feature=coverage type=model
//ff:what Represents a <sourcefile> element of a JaCoCo XML export holding per-line coverage
package coverage

// jacocoSourceFile is a <sourcefile> element holding per-line instruction and
// branch counters.
type jacocoSourceFile struct {
	Name  string       `xml:"name,attr"`
	Lines []jacocoLine `xml:"line"`
}
