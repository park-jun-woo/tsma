//ff:func feature=coverage type=implementation control=iteration dimension=1
//ff:what Parses a Go coverage profile file into a list of cover blocks
package coverage

import (
	"bufio"
	"os"
	"strings"
)

// parseCoverProfile parses a Go coverage profile file.
func parseCoverProfile(path string) ([]coverBlock, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var blocks []coverBlock
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}
		b, err := parseCoverLine(line)
		if err != nil {
			continue
		}
		blocks = append(blocks, b)
	}
	return blocks, scanner.Err()
}
