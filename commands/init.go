package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"mpkgr/format"
	"mpkgr/utils"
	"net/http"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v2"
)

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

func initPackage() format.Package {
	src := make(map[string]string)
	name := ""
	if len(os.Args) > 2 {
		req, err := http.NewRequest("GET", os.Args[2], nil)
		if err != nil {
			log.Fatal("Unable to download "+os.Args[2]+": ", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal("Unable to download "+os.Args[2]+": ", err)
		}
		defer resp.Body.Close()

		temp := strings.Split(os.Args[2], "/")
		filename := temp[len(temp)-1]

		f, err := os.Create(filename)
		if err != nil {
			log.Fatal("Unable to create "+os.Args[2]+": ", err)
		}
		defer f.Close()

		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"Downloading",
		)
		hash := sha256.New()
		if _, err := io.Copy(io.MultiWriter(f, bar, hash), resp.Body); err != nil {
			log.Fatal("Unable to create "+os.Args[2]+": ", err)
		}

		sum := hash.Sum(nil)
		sumString := hex.EncodeToString(sum)
		src[os.Args[2]] = sumString
		temp = strings.Split(filename, "-")
		name = temp[0]
		/*version := *semver.Version{}
		version, err = semver.NewVersion(temp[1])
		if err != nil {
			log.Println("Unable parse version, defaulting to 1.0")
			version, err = semver.NewVersion("1.0")
		}*/
	} else {
		log.Fatal("Must pass a url to a file")
	}
	pkg := format.Package{
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
