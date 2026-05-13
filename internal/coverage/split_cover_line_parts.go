//ff:func feature=coverage type=helper control=sequence
//ff:what Splits a coverage line into file/position, statement count, and hit count parts
package coverage

import (
	"fmt"
	"strings"
)

// splitCoverLineParts splits a coverage line into fileAndPos, stmtsStr, and countStr.
func splitCoverLineParts(line string) (string, string, string, error) {
	lastSpace := strings.LastIndex(line, " ")
	if lastSpace < 0 {
		return "", "", "", fmt.Errorf("invalid line: %s", line)
	}
	secondLastSpace := strings.LastIndex(line[:lastSpace], " ")
	if secondLastSpace < 0 {
		return "", "", "", fmt.Errorf("invalid line: %s", line)
	}

	return line[:secondLastSpace], line[secondLastSpace+1 : lastSpace], line[lastSpace+1:], nil
}
