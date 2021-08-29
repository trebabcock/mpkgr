package commands

import (
	"os"

	"mpkgr/format"
)

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
