//ff:func feature=cli type=helper control=sequence
//ff:what Prints the terminal banner shown when no TODO functions remain
package cli

import "fmt"

// printAllComplete prints the "all complete" banner shared by both next modes.
func printAllComplete() {
	fmt.Println("All functions complete!")
}
