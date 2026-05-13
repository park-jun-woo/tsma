//ff:func feature=coverage type=implementation control=sequence
//ff:what Parses a single line from a Go coverage profile into a coverBlock struct
package coverage

import (
	"strconv"
	"strings"
)

// parseCoverLine parses a single line from a Go coverage profile.
// Format: "file.go:startLine.startCol,endLine.endCol stmts count"
func parseCoverLine(line string) (coverBlock, error) {
	var b coverBlock

	fileAndPos, stmtsStr, countStr, err := splitCoverLineParts(line)
	if err != nil {
		return b, err
	}

	b.file, err = extractCoverFile(fileAndPos)
	if err != nil {
		return b, err
	}

	positions := fileAndPos[strings.LastIndex(fileAndPos, ":")+1:]
	if err := parseCoverPositions(positions, &b); err != nil {
		return b, err
	}

	b.stmts, _ = strconv.Atoi(stmtsStr)
	b.count, _ = strconv.Atoi(countStr)

	return b, nil
}
