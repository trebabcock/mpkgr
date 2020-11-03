package utils

import (
	"io"
	"log"
	"os"
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
