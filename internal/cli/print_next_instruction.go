//ff:func feature=cli type=helper control=sequence
//ff:what Prints loop-continuation instructions for LLM agents
package cli

import "fmt"

// printNextInstruction prints the instruction to run tsma next after completing the current task.
func printNextInstruction() {
	fmt.Println("  ▶ After completing the test, run `tsma next`. Don't stop. Keep going until complete.")
}

// printContinueInstruction prints the instruction to run tsma next for the next function.
func printContinueInstruction() {
	fmt.Println("  ▶ Run `tsma next` for the next function. Don't stop. Keep going until complete.")
}
