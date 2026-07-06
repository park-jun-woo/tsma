package smell

import (
	"os"
	"path/filepath"
	"testing"
)

// fakeTreeSitter writes an executable shell script that ignores its arguments
// and cats the canned XML to stdout, then points TSMA_TREE_SITTER at it. This
// exercises the full Scan* pipeline (ResolveCommand → Run → ParseXML →
// detectors) with no real tree-sitter CLI, mirroring the canned-XML style of
// cs_detectors_branch_test.go at the subprocess boundary.
func fakeTreeSitter(t *testing.T, xmlOut string) {
	t.Helper()
	dir := t.TempDir()
	xmlPath := filepath.Join(dir, "out.xml")
	if err := os.WriteFile(xmlPath, []byte(xmlOut), 0o644); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(dir, "fake-tree-sitter")
	body := "#!/bin/sh\nexec cat \"" + xmlPath + "\"\n"
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TSMA_TREE_SITTER", script)
}

// chdirToRemovedDir moves the process into a directory and deletes it, so
// os.Getwd (and therefore filepath.Abs on a relative path) fails. The original
// working directory is restored on cleanup.
func chdirToRemovedDir(t *testing.T) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), "gone")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("restore cwd: %v", err)
		}
	})
	if err := os.Remove(dir); err != nil {
		t.Fatal(err)
	}
}

// scanCsSuccessXML: one t.GetMethod() invocation -> exactly one TS-REFL-CS-001.
const scanCsSuccessXML = `<?xml version="1.0"?>
<sources>
  <source name="/x/T.cs">
    <compilation_unit srow="0" scol="0" erow="3" ecol="0">
      <invocation_expression srow="1" scol="0" erow="1" ecol="20">
        <member_access_expression field="function" srow="1" scol="0" erow="1" ecol="12">
          <identifier srow="1" scol="0" erow="1" ecol="1">t</identifier>
          <identifier field="name" srow="1" scol="2" erow="1" ecol="11">GetMethod</identifier>
        </member_access_expression>
        <argument_list field="arguments" srow="1" scol="12" erow="1" ecol="20"></argument_list>
      </invocation_expression>
    </compilation_unit>
  </source>
</sources>`

// scanJavaSuccessXML: one c.getDeclaredMethod() call -> exactly one TS-REFL-JV-001.
const scanJavaSuccessXML = `<?xml version="1.0"?>
<sources>
  <source name="/x/T.java">
    <program srow="0" scol="0" erow="3" ecol="0">
      <method_invocation srow="1" scol="0" erow="1" ecol="30">
        <identifier field="object" srow="1" scol="0" erow="1" ecol="1">c</identifier>
        <identifier field="name" srow="1" scol="2" erow="1" ecol="19">getDeclaredMethod</identifier>
        <argument_list field="arguments" srow="1" scol="19" erow="1" ecol="30"></argument_list>
      </method_invocation>
    </program>
  </source>
</sources>`

// scanRsSuccessXML: a top-level #[test] fn containing an unsafe block -> one
// test scope (rsTestScopeNodes) and exactly one TS-REFL-RS-001.
const scanRsSuccessXML = `<?xml version="1.0"?>
<sources>
  <source name="/x/lib.rs">
    <source_file srow="0" scol="0" erow="6" ecol="0">
      <attribute_item srow="0" scol="0" erow="0" ecol="7">
        <attribute srow="0" scol="2" erow="0" ecol="6">
          <identifier srow="0" scol="2" erow="0" ecol="6">test</identifier>
        </attribute>
      </attribute_item>
      <function_item srow="1" scol="0" erow="4" ecol="1">
        <identifier field="name" srow="1" scol="3" erow="1" ecol="8">works</identifier>
        <block field="body" srow="1" scol="11" erow="4" ecol="1">
          <unsafe_block srow="2" scol="4" erow="3" ecol="5">
            <block srow="2" scol="11" erow="3" ecol="5"></block>
          </unsafe_block>
        </block>
      </function_item>
    </source_file>
  </source>
</sources>`

// scanTSSuccessXML: one Reflect.get member expression -> exactly one TS-REFL-TS-002.
const scanTSSuccessXML = `<?xml version="1.0"?>
<sources>
  <source name="/x/a.test.ts">
    <program srow="0" scol="0" erow="2" ecol="0">
      <member_expression srow="1" scol="0" erow="1" ecol="11">
        <identifier field="object" srow="1" scol="0" erow="1" ecol="7">Reflect</identifier>
        <property_identifier field="property" srow="1" scol="8" erow="1" ecol="11">get</property_identifier>
      </member_expression>
    </program>
  </source>
</sources>`

// truncated XML: Run succeeds (non-empty stdout, exit 0) but ParseXML fails,
// covering the ParseFile-error branch of every Scan*.
const scanBadXML = `<sources><source`

// scanFakeCLICases drives the shared branch matrix over the four tree-sitter
// scanners: success (detectors fire), ParseFile error (malformed XML), and
// filepath.Abs error (deleted cwd + relative path).
var scanFakeCLICases = []struct {
	name       string
	scan       func(string) ([]Finding, error)
	path       string
	successXML string
	wantRule   string
}{
	{"cs", ScanCs, "T.cs", scanCsSuccessXML, "TS-REFL-CS-001"},
	{"java", ScanJava, "T.java", scanJavaSuccessXML, "TS-REFL-JV-001"},
	{"rs", ScanRs, "lib.rs", scanRsSuccessXML, "TS-REFL-RS-001"},
	{"ts", ScanTS, "a.test.ts", scanTSSuccessXML, "TS-REFL-TS-002"},
}

// TestScanFakeCLISuccess covers each scanner's happy path end to end with a
// fake CLI emitting canned XML: parse succeeds and the detectors fire.
func TestScanFakeCLISuccess(t *testing.T) {
	for _, tc := range scanFakeCLICases {
		t.Run(tc.name, func(t *testing.T) {
			fakeTreeSitter(t, tc.successXML)
			findings, err := tc.scan(tc.path)
			if err != nil {
				t.Fatalf("scan: %v", err)
			}
			if len(findings) != 1 || findings[0].Rule != tc.wantRule {
				t.Errorf("findings = %+v, want exactly one %s", findings, tc.wantRule)
			}
			if findings[0].File != tc.path {
				t.Errorf("finding file = %q, want %q", findings[0].File, tc.path)
			}
		})
	}
}

// TestScanFakeCLIParseError covers each scanner's ParseFile-error branch: the
// fake CLI emits truncated XML, so ParseXML fails and (nil, err) is returned.
func TestScanFakeCLIParseError(t *testing.T) {
	for _, tc := range scanFakeCLICases {
		t.Run(tc.name, func(t *testing.T) {
			fakeTreeSitter(t, scanBadXML)
			findings, err := tc.scan(tc.path)
			if err == nil {
				t.Error("expected parse error on truncated XML")
			}
			if findings != nil {
				t.Errorf("expected nil findings, got %+v", findings)
			}
		})
	}
}

// TestScanCsDirectFakeCLI drives ScanCs by name through the fake-CLI matrix —
// success, ParseFile error (truncated XML), and the filepath.Abs-error fallback
// (deleted cwd) — so the branches attribute to ScanCs itself.
func TestScanCsDirectFakeCLI(t *testing.T) {
	fakeTreeSitter(t, scanCsSuccessXML)
	findings, err := ScanCs("T.cs")
	if err != nil || len(findings) != 1 || findings[0].Rule != "TS-REFL-CS-001" {
		t.Fatalf("success = (%+v, %v), want one TS-REFL-CS-001", findings, err)
	}
	fakeTreeSitter(t, scanBadXML)
	if got, err := ScanCs("T.cs"); err == nil || got != nil {
		t.Errorf("parse error = (%+v, %v), want (nil, err)", got, err)
	}
	fakeTreeSitter(t, scanCsSuccessXML)
	chdirToRemovedDir(t)
	if got, err := ScanCs("T.cs"); err != nil || len(got) != 1 {
		t.Errorf("abs fallback = (%+v, %v), want one finding", got, err)
	}
}

// TestScanJavaDirectFakeCLI mirrors TestScanCsDirectFakeCLI for ScanJava.
func TestScanJavaDirectFakeCLI(t *testing.T) {
	fakeTreeSitter(t, scanJavaSuccessXML)
	findings, err := ScanJava("T.java")
	if err != nil || len(findings) != 1 || findings[0].Rule != "TS-REFL-JV-001" {
		t.Fatalf("success = (%+v, %v), want one TS-REFL-JV-001", findings, err)
	}
	fakeTreeSitter(t, scanBadXML)
	if got, err := ScanJava("T.java"); err == nil || got != nil {
		t.Errorf("parse error = (%+v, %v), want (nil, err)", got, err)
	}
	fakeTreeSitter(t, scanJavaSuccessXML)
	chdirToRemovedDir(t)
	if got, err := ScanJava("T.java"); err != nil || len(got) != 1 {
		t.Errorf("abs fallback = (%+v, %v), want one finding", got, err)
	}
}

// TestScanRsDirectFakeCLI mirrors TestScanCsDirectFakeCLI for ScanRs; the
// success XML carries one #[test] fn scope so the scope loop and detectors run.
func TestScanRsDirectFakeCLI(t *testing.T) {
	fakeTreeSitter(t, scanRsSuccessXML)
	findings, err := ScanRs("lib.rs")
	if err != nil || len(findings) != 1 || findings[0].Rule != "TS-REFL-RS-001" {
		t.Fatalf("success = (%+v, %v), want one TS-REFL-RS-001", findings, err)
	}
	fakeTreeSitter(t, scanBadXML)
	if got, err := ScanRs("lib.rs"); err == nil || got != nil {
		t.Errorf("parse error = (%+v, %v), want (nil, err)", got, err)
	}
	fakeTreeSitter(t, scanRsSuccessXML)
	chdirToRemovedDir(t)
	if got, err := ScanRs("lib.rs"); err != nil || len(got) != 1 {
		t.Errorf("abs fallback = (%+v, %v), want one finding", got, err)
	}
}

// TestScanTSDirectFakeCLI mirrors TestScanCsDirectFakeCLI for ScanTS.
func TestScanTSDirectFakeCLI(t *testing.T) {
	fakeTreeSitter(t, scanTSSuccessXML)
	findings, err := ScanTS("a.test.ts")
	if err != nil || len(findings) != 1 || findings[0].Rule != "TS-REFL-TS-002" {
		t.Fatalf("success = (%+v, %v), want one TS-REFL-TS-002", findings, err)
	}
	fakeTreeSitter(t, scanBadXML)
	if got, err := ScanTS("a.test.ts"); err == nil || got != nil {
		t.Errorf("parse error = (%+v, %v), want (nil, err)", got, err)
	}
	fakeTreeSitter(t, scanTSSuccessXML)
	chdirToRemovedDir(t)
	if got, err := ScanTS("a.test.ts"); err != nil || len(got) != 1 {
		t.Errorf("abs fallback = (%+v, %v), want one finding", got, err)
	}
}

// TestScanFakeCLIAbsError covers each scanner's filepath.Abs-error branch: with
// the working directory deleted, Abs on a relative path fails and the scan
// falls back to the relative path verbatim (the pipeline still runs — the fake
// CLI cats an absolute XML file, and ParseFile matches the sole source).
func TestScanFakeCLIAbsError(t *testing.T) {
	for _, tc := range scanFakeCLICases {
		t.Run(tc.name, func(t *testing.T) {
			fakeTreeSitter(t, tc.successXML)
			chdirToRemovedDir(t)
			findings, err := tc.scan(tc.path)
			if err != nil {
				t.Fatalf("scan: %v", err)
			}
			if len(findings) != 1 || findings[0].Rule != tc.wantRule {
				t.Errorf("findings = %+v, want exactly one %s", findings, tc.wantRule)
			}
		})
	}
}
