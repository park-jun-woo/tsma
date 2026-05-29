//ff:func feature=match type=helper control=sequence
//ff:what Maps a Java source directory to its corresponding test source directory
package match

import "strings"

// javaTestDir maps the directory of a Java source file to the directory where
// its test would live. For the standard Maven/Gradle layout it rewrites the
// first "src/main/" segment to "src/test/". When the source is not under a
// main source tree, the same directory is returned (tests beside sources).
func javaTestDir(srcDir string) string {
	slashed := strings.ReplaceAll(srcDir, "\\", "/")
	if strings.Contains(slashed, "src/main/") {
		return strings.Replace(slashed, "src/main/", "src/test/", 1)
	}
	return slashed
}
