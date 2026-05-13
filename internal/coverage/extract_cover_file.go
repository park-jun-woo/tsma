//ff:func feature=coverage type=helper control=sequence
//ff:what Extracts the file path from the position part of a Go coverage line
package coverage

import (
	"fmt"
	"strings"
)

// extractCoverFile extracts the file path from the position part of a coverage line.
func extractCoverFile(positionPart string) (string, error) {
	colonIdx := strings.LastIndex(positionPart, ":")
	if colonIdx < 0 {
		return "", fmt.Errorf("invalid line: no colon in %s", positionPart)
	}

	commaIdx := strings.Index(positionPart[colonIdx:], ",")
	if commaIdx < 0 {
		return "", fmt.Errorf("invalid line: no comma in %s", positionPart)
	}

	return positionPart[:colonIdx], nil
}
