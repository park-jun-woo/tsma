//ff:type feature=model type=model
//ff:what Represents a single function to be tested
package model

type Function struct {
	QualifiedName string `json:"qualified_name"`
	Name          string `json:"name"`
	File          string `json:"file"`
	StartLine     int    `json:"start_line"`
	EndLine       int    `json:"end_line"`
	Exported      bool   `json:"exported"`
	TestFile      string `json:"test_file,omitempty"`
	Status        string `json:"status"`
	FailOutput    string `json:"fail_output,omitempty"`
}
