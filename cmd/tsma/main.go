//ff:func feature=cli type=command control=sequence
//ff:what Entry point that delegates to the CLI Execute function
package main

import "github.com/park-jun-woo/tsma/internal/cli"

// Version is set at build time via ldflags.
var Version = "dev"

func main() {
	cli.Execute(Version)
}
