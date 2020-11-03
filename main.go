package main

import (
	"fmt"
	"os"

	"mpkgr/commands"

	"github.com/mingrammer/cfmt"
)

func main() {
	commands.SetCommands()

	if val, ok := commands.Commands[os.Args[1]]; ok {
		val.Func()
	} else {
		cfmt.Errorln("Unknown command: " + os.Args[1])
		fmt.Println("Run 'mpkgr help' for usage.")
	}
}
