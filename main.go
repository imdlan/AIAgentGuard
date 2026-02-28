package main

import (
	"github.com/imdlan/AIAgentGuard/cmd"
)

	// version is set at build time using ldflags
var version = "v1.2.0-beta"

func main() {
	cmd.Execute()
}
