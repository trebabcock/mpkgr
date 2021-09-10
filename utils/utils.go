package utils

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Assert compares two values and returns an error if they are not equal
func Assert(assertion bool, exception string) {
	if !assertion {
		log.Fatal(exception)
	}
}

// WriteStringToNewFile writes a string to a new file
func WriteStringToNewFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}

	return file.Sync()
}

// TarDir creates a tarball from a directory
func TarDir(sourcedir string, destinationfile string) {
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
