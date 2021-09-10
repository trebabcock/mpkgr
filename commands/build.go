package commands

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"mpkgr/format"
	"mpkgr/pkg"

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

	fileinfo, err := ioutil.ReadDir("./build_output")
	if err != nil {
		log.Fatal("Unable to read files: ", err)
	}

	files := []*os.File{}
	for _, f := range fileinfo {
		file, err := os.Open("./build_output/" + f.Name())
		if err != nil {
			log.Fatal("Unable to open ", f.Name(), ": ", err)
		}
		files = append(files, file)
	}

	fileBytes := []byte{}

	for _, f := range files {
		temp, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal("Unable to read ", f.Name(), ": ", err)
		}
		fileBytes = append(fileBytes, temp)
	}

	tarDir("./build_output", "content.tar")

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

func tarDir(sourcedir string, destinationfile string) {
	dir, err := os.Open(sourcedir)
	checkerror(err)
	defer dir.Close()

	// get list of files
	files, err := dir.Readdir(0)
	checkerror(err)

	// create tar file
	tarfile, err := os.Create(destinationfile)
	checkerror(err)
	defer tarfile.Close()

	var fileWriter io.WriteCloser = tarfile

	tarfileWriter := tar.NewWriter(fileWriter)
	defer tarfileWriter.Close()

	for _, fileInfo := range files {

		if fileInfo.IsDir() {
			continue
		}

		file, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())
		checkerror(err)
		defer file.Close()

		// prepare the tar header
		header := new(tar.Header)
		header.Name = file.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = tarfileWriter.WriteHeader(header)
		checkerror(err)

		_, err = io.Copy(tarfileWriter, file)
		checkerror(err)
	}
}

func checkerror(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
