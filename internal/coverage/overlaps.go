//ff:func feature=coverage type=helper control=sequence
//ff:what Checks if a cover block overlaps with a target function file and line range
package coverage

import "strings"

func overlaps(blockFile, targetFile string, blockStart, blockEnd, funcStart, funcEnd int) bool {
	if !strings.HasSuffix(blockFile, targetFile) {
		return false
	}
	return blockStart >= funcStart && blockStart <= funcEnd
}
