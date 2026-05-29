//ff:func feature=cli type=helper control=sequence
//ff:what Prints an instruction to rename a misnamed test file to its canonical name
package cli

import "fmt"

// printRenameInstruction prints the instruction to rename a misnamed test file
// to its canonical name, then run tsma next.
func printRenameInstruction(misnamed, canonical string) {
	fmt.Printf("  ▶ Test file misnamed: rename `%s` → `%s`, then run `tsma next`.\n", misnamed, canonical)
}
