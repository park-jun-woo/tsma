//ff:func feature=endpoint type=helper control=sequence
//ff:what Resolves pending decorators against a def line or clears them on non-decorator lines
package endpoint

import "strings"

func resolvePendingDecs(pendingDec []pendingDecorator, lines []string, line string, idx, lineNum int, relPath string) ([]pyRoute, []pendingDecorator) {
	if len(pendingDec) == 0 {
		return nil, pendingDec
	}

	if dm := pyDefRe.FindStringSubmatch(line); dm != nil {
		routes := buildRoutesFromDecs(pendingDec, dm, lines, idx, lineNum, relPath)
		return routes, nil
	}

	trimmed := strings.TrimSpace(line)
	if trimmed != "" && !strings.HasPrefix(trimmed, "@") && !strings.HasPrefix(trimmed, "#") {
		return nil, nil
	}

	return nil, pendingDec
}
