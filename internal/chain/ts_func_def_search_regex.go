//ff:func feature=chain type=helper control=sequence
//ff:what Builds a regex to match TS/JS function definitions for a specific name
package chain

import "regexp"

// tsFuncDefSearchRegex builds a regex to match function definitions for a specific name.
func tsFuncDefSearchRegex(name string) *regexp.Regexp {
	escaped := regexp.QuoteMeta(name)
	return regexp.MustCompile(
		`(?:export\s+)?(?:async\s+)?function\s+` + escaped + `\s*\(` +
			`|(?:export\s+)?(?:const|let|var)\s+` + escaped + `\s*=` +
			`|` + escaped + `\s*\([^)]*\)\s*[:{]`,
	)
}
