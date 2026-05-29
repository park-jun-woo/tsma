//ff:type feature=coverage type=model
//ff:what Represents a <line> element of a JaCoCo sourcefile with instruction and branch counters
package coverage

// jacocoLine is a <line> element. The attributes are JaCoCo's per-line
// counters: mi/ci are missed/covered instructions, mb/cb are missed/covered
// branches, and nr is the 1-based source line number.
type jacocoLine struct {
	Nr int `xml:"nr,attr"`
	Mi int `xml:"mi,attr"`
	Ci int `xml:"ci,attr"`
	Mb int `xml:"mb,attr"`
	Cb int `xml:"cb,attr"`
}
