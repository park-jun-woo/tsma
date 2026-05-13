//ff:func feature=runner type=helper control=sequence
//ff:what Reads and returns raw bytes from a file path
package runner

import "os"

func readFileBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}
