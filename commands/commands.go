package commands

import (
	"fmt"

	"github.com/gosuri/uitable"
)

// Command has a name, description, and function for each command
type Command struct {
	Name string
	Desc string
	Run  func()
}

// Commands maps a command string to a Command object
var Commands = make(map[string]Command)

// SetCommands sets the commands for mpkgr
func SetCommands() {

	help := Command{
		Name: "help",
		Desc: "Show this help message",
		Run:  usage,
	}

	init := Command{
		Name: "init",
		Desc: "Generate package.yml template",
		Run:  initialize,
	}

	build := Command{
		Name: "build",
		Desc: "Build a package",
		Run:  BuildPackage,
	}

	info := Command{
		Name: "info",
		Desc: "Display info about a package",
		Run:  info,
	}

	Commands["help"] = help
	Commands["init"] = init
	Commands["build"] = build
	Commands["info"] = info
}

func usage() {
	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	table.AddRow("COMMAND", "DESCRIPTION")
	table.AddRow("")
	for _, v := range Commands {
		table.AddRow(v.Name, v.Desc)
	}

	fmt.Println(table)
}

func info() {

}
