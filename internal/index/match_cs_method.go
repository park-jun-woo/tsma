//ff:func feature=index type=helper control=sequence lang=csharp
//ff:what Matches a C# method or constructor declaration line and returns its name
package index

// matchCsMethod reports whether the trimmed line declares a method or
// constructor whose body opens on the same line, returning the declared name.
//
// It rejects type declarations (class/struct/interface/record/enum) and
// control-flow statements (if/for/foreach/...) that share the `name(...) {`
// shape. A regular method requires a return-type/modifier prefix; a constructor
// is recognised by an uppercase-led name with optional access modifiers
// (best-effort).
func matchCsMethod(trimmed string) (string, bool) {
	if csTypePattern.MatchString(trimmed) {
		return "", false
	}
	if m := csMethodPattern.FindStringSubmatch(trimmed); m != nil {
		name := m[1]
		if !isCsControlKeyword(name) {
			return name, true
		}
	}
	if m := csConstructorPattern.FindStringSubmatch(trimmed); m != nil {
		name := m[1]
		if !isCsControlKeyword(name) {
			return name, true
		}
	}
	return "", false
}
