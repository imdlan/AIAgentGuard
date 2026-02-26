package main

import (
	"github.com/imdlan/AIAgentGuard/cmd"
)

	// version is set at build time using ldflags
var version = "dev"

func main() {
	cmd.Execute()
}
