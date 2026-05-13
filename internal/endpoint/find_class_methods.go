//ff:func feature=endpoint type=implementation control=iteration dimension=1
//ff:what Finds HTTP method definitions inside a Django class-based view body
package endpoint

import "strings"

// findClassMethods finds HTTP method definitions (get, post, etc.) inside a class body.
func findClassMethods(lines []string, classIdx int, classIndent string) []classMethod {
	classIndentLen := effectiveIndent(classIndent)
	var methods []classMethod

	for i := classIdx + 1; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		lineIndentLen := effectiveIndentStr(line)
		if lineIndentLen <= classIndentLen && trimmed != "" {
			break
		}

		if m := djangoClassMethodRe.FindStringSubmatch(line); m != nil {
			methodIndent := leadingWhitespace(line)
			endLine := findPyFuncEndDjango(lines, i, methodIndent)
			methods = append(methods, classMethod{
				name:      m[1],
				startLine: i + 1,
				endLine:   endLine,
			})
		}
	}

	return methods
}
