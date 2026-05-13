//ff:func feature=coverage type=helper control=sequence
//ff:what Parses startLine.startCol,endLine.endCol into a coverBlock
package coverage

import (
	"fmt"
	"strconv"
	"strings"
)

// parseCoverPositions parses "startLine.startCol,endLine.endCol" into a coverBlock.
func parseCoverPositions(positions string, b *coverBlock) error {
	parts := strings.Split(positions, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid positions: %s", positions)
	}

	start := strings.Split(parts[0], ".")
	end := strings.Split(parts[1], ".")
	if len(start) != 2 || len(end) != 2 {
		return fmt.Errorf("invalid position format: %s", positions)
	}

	b.startLine, _ = strconv.Atoi(start[0])
	b.startCol, _ = strconv.Atoi(start[1])
	b.endLine, _ = strconv.Atoi(end[0])
	b.endCol, _ = strconv.Atoi(end[1])

	return nil
}
