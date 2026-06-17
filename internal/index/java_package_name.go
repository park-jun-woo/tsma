//ff:func feature=index type=helper control=sequence lang=java
//ff:what javaPackageName: extracts the dotted package name from a parsed Java program by finding the package_declaration child and joining its identifier leaf texts with "." (so a nested scoped_identifier com→example→app yields "com.example.app"). Returns "" for a default-package file. Feeds buildJavaQualifiedName exactly as the line-based javaPackagePattern does.
package index

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// javaPackageName returns the dotted package name declared in the file, or "".
func javaPackageName(root *treesitter.Node) string {
	pkgDecl := root.ChildByType("package_declaration")
	if pkgDecl == nil {
		return ""
	}
	var parts []string
	treesitter.Walk(pkgDecl, func(n *treesitter.Node) bool {
		if n.Type == "identifier" && n.Text != "" {
			parts = append(parts, n.Text)
		}
		return true
	})
	return strings.Join(parts, ".")
}
