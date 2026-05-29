//ff:func feature=cli type=helper control=sequence
//ff:what Prints continue instruction for LLM agents to run tsma next
package cli

import "fmt"

// printContinueInstruction prints the instruction to run tsma next for the next function.
func printContinueInstruction() {
	fmt.Println("  ▶ Run `tsma next` for the next function. Don't stop. Keep going until complete.")
}
