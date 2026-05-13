//ff:func feature=chain type=implementation control=iteration dimension=1
//ff:what Parses a single Python file and extracts all function definitions with metadata
package chain

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// parsePyFunctions parses a single Python file and extracts function definitions.
func parsePyFunctions(filePath, relPath string) map[string]*pyFuncInfo {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil
	}

	funcs := make(map[string]*pyFuncInfo)
	modKey := strings.TrimSuffix(relPath, ".py")
	modKey = strings.ReplaceAll(modKey, string(filepath.Separator), ".")

	var currentClass string
	var classIndentLen int

	for i, line := range lines {
		if m := pyClassDefRe.FindStringSubmatch(line); m != nil {
			currentClass = m[2]
			classIndentLen = pyEffectiveIndent(m[1])
			continue
		}

		m := pyFuncDefRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		fi := buildPyFuncInfo(m, lines, i, relPath)
		key := buildPyFuncKey(modKey, m[2], pyEffectiveIndent(m[1]), currentClass, classIndentLen)
		if pyEffectiveIndent(m[1]) <= classIndentLen {
			currentClass = ""
		}
		funcs[key] = fi
	}

	return funcs
}
