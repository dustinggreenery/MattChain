package main

import (
	"os"

	"github.com/dustinggreenery/MattChain/cli"
)

func main() {
	defer os.Exit(0);
	cmd := cli.CommandLine{};
	cmd.Run();
}