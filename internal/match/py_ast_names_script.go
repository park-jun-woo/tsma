//ff:func feature=match type=helper control=sequence lang=python
//ff:what pyAstNamesScript: the embedded Python program (run via `python -c`) that ast-parses one test file and dumps the set of names it references — every Call target (bare name, or the trailing attribute of obj.method()) plus every imported name (from-import members and the top package of plain imports). It is the Python analogue of collectTSCalledNames: the source function's bare Name is looked up against this set. Test-harness noise (pytest, assert helpers) is collected too but harmlessly. Exit 2 on SyntaxError so the Go side skips the file.
package match

// pyAstNamesScript is passed to `python -c`; argv[1] is the test file. It prints
// a JSON array of referenced names (called + imported).
const pyAstNamesScript = `import ast, json, sys

src = open(sys.argv[1]).read()
try:
    tree = ast.parse(src)
except SyntaxError:
    sys.exit(2)
names = set()
for node in ast.walk(tree):
    if isinstance(node, ast.Call):
        f = node.func
        if isinstance(f, ast.Name):
            names.add(f.id)
        elif isinstance(f, ast.Attribute):
            names.add(f.attr)
    elif isinstance(node, ast.ImportFrom):
        for a in node.names:
            names.add(a.asname or a.name)
    elif isinstance(node, ast.Import):
        for a in node.names:
            names.add((a.asname or a.name).split(".")[0])
json.dump(sorted(names), sys.stdout)
`
