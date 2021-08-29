package main

import (
	"fmt"
	"os"

	"mpkgr/commands"

	"github.com/mingrammer/cfmt"
)

func main() {
	commands.SetCommands()

	if len(os.Args) < 2 {
		cfmt.Errorln("Missing command")
		fmt.Println("Run 'mpkgr help' for usage.")
		return
	}

	if val, ok := commands.Commands[os.Args[1]]; ok {
		val.Run()
	} else {
		cfmt.Errorln("Unknown command: " + os.Args[1])
		fmt.Println("Run 'mpkgr help' for usage.")
	}
}
