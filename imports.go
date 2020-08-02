package interfaces

import (
	"io/ioutil"

	"golang.org/x/tools/imports"
)

// FormatFile runs goimports on given filename.
func FormatFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	src, err := imports.Process(filename, buf, nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, src, 0644)
}
