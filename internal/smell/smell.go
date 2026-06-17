//ff:type feature=smell type=model
//ff:what test-smell 탐지의 공용 타입. Finding은 한 escape-hatch 위반(Rule ID·파일·라인·사람용 Note)을 담는다. smell은 커버리지와 직교하는 "테스트가 내부를 강제로 비트는가" 축으로, ScanGo(scan_go.go)가 한 _test.go를 go/ast로 파싱해 []Finding을 낸다. 정당한 reflect.DeepEqual/TypeOf/ValueOf는 절대 Finding이 아니다(위양성 0).

package smell

// Finding is a single escape-hatch detection in a test file: the rule it
// violates, the location (file + line), and a one-line human note. It is the
// language-agnostic carrier so a tsmagate measurement can hold findings from any
// ScanX detector (Go first; ScanTS/ScanPy later, rulebook §7).
type Finding struct {
	// Rule is the smell rule ID, e.g. "TS-REFL-001".
	Rule string
	// File is the path passed to the scanner (root-relative or absolute as given).
	File string
	// Line is the 1-based source line of the offending node.
	Line int
	// Note is a short human-readable description of what was found.
	Note string
}
