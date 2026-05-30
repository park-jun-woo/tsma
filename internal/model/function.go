//ff:type feature=model type=model
//ff:what Represents a single function to be tested
package model

type Function struct {
	QualifiedName string   `json:"qualified_name"`
	Name          string   `json:"name"`
	Receiver      string   `json:"receiver,omitempty"`
	File          string   `json:"file"`
	StartLine     int      `json:"start_line"`
	EndLine       int      `json:"end_line"`
	Exported      bool     `json:"exported"`
	TestFile      string   `json:"test_file,omitempty"`
	TestFiles     []string `json:"test_files,omitempty"`
	TestMtime     string   `json:"test_mtime,omitempty"`
	Status        string   `json:"status"`
	CoveragePct   float64  `json:"coverage_pct,omitempty"`
	Attempt       int      `json:"attempt,omitempty"`
	FailOutput    string   `json:"fail_output,omitempty"`
}
