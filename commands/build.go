package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"mpkgr/format"
	"mpkgr/pkg"
	"mpkgr/utils"

	"github.com/kelindar/binary"
	"gopkg.in/yaml.v2"
)

// BuildPackage builds the package
func BuildPackage() {

	pkgConfig := format.Package{}

	configFile, err := os.Open("package.yml")
	if err != nil {
		log.Fatal("Unable to open package.yml: ", err)
	}

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal("Unable to read package.yml: ", err)
	}

	if err = yaml.Unmarshal(configBytes, &pkgConfig); err != nil {
		log.Fatal("Unable to parse package.yml: ", err)
	}

	buildSetup(&pkgConfig)

	header := pkg.PackageHeader{
		MpkgVersion: "0.0.1",
	}

	meta := pkg.PackageMeta{
		Name:         pkgConfig.Name,
		Version:      pkgConfig.Version,
		Dependencies: pkgConfig.Rundeps,
		Scripts:      pkgConfig.Install,
	}

	utils.TarDir("./build_output", "content.tar")

	tarFile, err := os.Open("./content.tar")
	if err != nil {
		log.Fatal("Unable to open tarball: ", err)
	}

	tarBytes, err := ioutil.ReadAll(tarFile)
	if err != nil {
		log.Fatal("Unable to read tarball: ", err)
	}

	content := pkg.PackageContent{Files: tarBytes}

	finalPackage := pkg.Package{
		Header:  header,
		Meta:    meta,
		Content: content,
	}

	encoded, err := binary.Marshal(finalPackage)
	if err != nil {
		log.Fatal("Unable to marshal binary data: ", err)
	}
	reader := bytes.NewReader(encoded)

	file, err := os.Create(fmt.Sprintf("%s-%s.mpkg", pkgConfig.Name, pkgConfig.Version))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err = io.Copy(file, reader); err != nil {
		log.Fatal("Unable to write mpkg file: ", err)
	}
}

func buildSetup(config *format.Package) {
	for _, c := range config.Setup {
		cmd := exec.Command("bash", "-c", c)
		stdout, err := cmd.Output()
		if err != nil {
			log.Fatal("Unable to run Setup commands: ", err)
			return
		}
		fmt.Println(string(stdout))
	}
}

func checkerror(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
