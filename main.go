package main

import (
	"github.com/imdlan/AIAgentGuard/cmd"
)

	// version is set at build time using ldflags
var version = "v1.1.0"

func main() {
	cmd.Execute()
}
