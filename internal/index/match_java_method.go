//ff:func feature=index type=helper control=sequence
//ff:what Matches a Java method or constructor declaration line and returns its name
package index

// matchJavaMethod reports whether the trimmed line declares a method or
// constructor whose body opens on the same line, returning the declared name.
//
// It rejects type declarations (class/interface/enum/record) and control-flow
// statements (if/for/while/...) that share the `name(...) {` shape. A regular
// method requires a return-type/modifier prefix; a constructor is recognised by
// an uppercase-led name with optional access modifiers (best-effort).
func matchJavaMethod(trimmed string) (string, bool) {
	if javaTypePattern.MatchString(trimmed) {
		return "", false
	}
	if m := javaMethodPattern.FindStringSubmatch(trimmed); m != nil {
		name := m[1]
		if !isJavaControlKeyword(name) {
			return name, true
		}
	}
	if m := javaConstructorPattern.FindStringSubmatch(trimmed); m != nil {
		name := m[1]
		if !isJavaControlKeyword(name) {
			return name, true
		}
	}
	return "", false
}
