package commands

import (
	"os"

	"mpkgr/format"
)

type Package struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Release     int               `yaml:"release"`
	Source      map[string]string `yaml:"source"`
	License     []string          `yaml:"license"`
	Summary     string            `yaml:"summary"`
	Component   string            `yaml:"component"`
	Description string            `yaml:"description"`
	Builddeps   []string          `yaml:"builddeps"`
	Rundeps     []string          `yaml:"rundeps"`
	Homepage    string            `yaml:"homepage"`
	Setup       []string          `yaml:"setup"`
	Build       []string          `yaml:"build"`
	Install     []string          `yaml:"install"`
}

// BuildPackage builds the package
func BuildPackage() {
	header := format.ConstructHeader(1)

	p := format.DefaultPayload()

	file, err := os.Create("test.mpkg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	header.Validate()
	header.Encode(file)
	p.Encode(file)
}
