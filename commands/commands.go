package commands

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"mpkgr/utils"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gosuri/uitable"
	"gopkg.in/yaml.v2"
)

// Command has a name, description, and function for each command
type Command struct {
	Name string
	Desc string
	Func func()
}

// Commands maps a command string to a Command object
var Commands = make(map[string]Command)

// SetCommands sets the commands for mpkgr
func SetCommands() {

	help := Command{
		Name: "help",
		Desc: "Show this help message",
		Func: usage,
	}

	init := Command{
		Name: "init",
		Desc: "Generate package.yml template",
		Func: initialize,
	}

	build := Command{
		Name: "build",
		Desc: "Build a package",
		Func: BuildPackage,
	}

	info := Command{
		Name: "info",
		Desc: "Display info about a package",
		Func: info,
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

func initialize() {
	pkg := initPackage()

	d, err := yaml.Marshal(&pkg)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.WriteStringToNewFile("package.yml", string(d))
	if err != nil {
		log.Fatal(err)
	}
}

func info() {

}

func initPackage() Package {
	src := make(map[string]string)
	name := ""
	if len(os.Args) > 2 {
		resp, err := http.Get(os.Args[2])
		//log.Fatal(resp, os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		out, err := os.Create(path.Base(resp.Request.URL.Path))
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		if _, err = io.Copy(out, resp.Body); err != nil {
			log.Fatal(err)
		}

		out.Seek(0, 0)

		hash := sha256.New()
		if _, err := io.Copy(hash, out); err != nil {
			log.Fatal(err)
		}

		sum := hash.Sum(nil)

		src[os.Args[2]] = string(sum)
		name = strings.Trim(path.Base(resp.Request.URL.Path), ".")
	} else {
		log.Fatal("Must pass a url to a file")
	}
	pkg := Package{
		Name:        name,
		Version:     "1.0",
		Release:     1,
		Source:      src,
		License:     []string{"MIT"},
		Summary:     "A short summary of the program.",
		Description: "A longer description of the program.",
		Builddeps:   []string{"EDIT THIS"},
		Rundeps:     []string{"EDIT THIS"},
		Homepage:    "https://www.example.com/",
		Setup:       []string{""},
		Build:       []string{"%make"},
		Install:     []string{"%make-install"},
	}

	return pkg
}
