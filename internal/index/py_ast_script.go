//ff:func feature=index type=helper control=sequence lang=python
//ff:what pyAstDefScript: the embedded Python program (run via `python -c`) that parses one .py file with the built-in `ast` module and dumps every FunctionDef/AsyncFunctionDef as JSON (name, start/end line, col, owning class receiver, decorators). It recurses with ast.iter_child_nodes tracking the enclosing class so a method gets receiver=class while a nested function inside a method gets receiver="" (Phase005b §1 parent-tracking memo). ast.parse only parses — it never imports/executes the module — so D1 is side-effect free. Exit 2 on SyntaxError so the Go side per-file-falls-back to the line indexer.
package index

// pyAstDefScript is passed to `python -c`; argv[1] is the target .py file. It
// prints a JSON array of {name,start_line,end_line,col,receiver,decorators}.
const pyAstDefScript = `import ast, json, sys

def main():
    src = open(sys.argv[1]).read()
    try:
        tree = ast.parse(src)
    except SyntaxError:
        sys.exit(2)
    out = []
    def visit(node, cls):
        for child in ast.iter_child_nodes(node):
            if isinstance(child, (ast.FunctionDef, ast.AsyncFunctionDef)):
                decos = []
                for d in child.decorator_list:
                    try:
                        decos.append(ast.unparse(d))
                    except Exception:
                        pass
                out.append({
                    "name": child.name,
                    "start_line": child.lineno,
                    "end_line": getattr(child, "end_lineno", child.lineno),
                    "col": child.col_offset,
                    "receiver": cls or "",
                    "decorators": decos,
                })
                visit(child, "")
            elif isinstance(child, ast.ClassDef):
                visit(child, child.name)
            else:
                visit(child, cls)
    visit(tree, "")
    json.dump(out, sys.stdout)

main()
`
