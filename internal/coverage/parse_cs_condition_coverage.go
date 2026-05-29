//ff:func feature=coverage type=helper control=sequence lang=csharp
//ff:what Parses a Cobertura condition-coverage attribute into covered and total branch counts
package coverage

import (
	"strconv"
	"strings"
)

// parseCsConditionCoverage extracts the covered and total branch counts from a
// Cobertura condition-coverage attribute of the form "50% (1/2)". The fraction
// inside the parentheses is parsed; (0, 0) is returned when the attribute is
// absent or malformed.
func parseCsConditionCoverage(cc string) (covered, total int) {
	open := strings.IndexByte(cc, '(')
	close := strings.IndexByte(cc, ')')
	if open < 0 || close < 0 || close <= open {
		return 0, 0
	}
	frac := cc[open+1 : close]
	slash := strings.IndexByte(frac, '/')
	if slash < 0 {
		return 0, 0
	}
	c, err1 := strconv.Atoi(strings.TrimSpace(frac[:slash]))
	tot, err2 := strconv.Atoi(strings.TrimSpace(frac[slash+1:]))
	if err1 != nil || err2 != nil {
		return 0, 0
	}
	return c, tot
}
