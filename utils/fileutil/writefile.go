package fileutil

import (
	"io/ioutil"
	"os"
	"path"
)

func OverWrite(filepath string, body []byte) error {
	dir := path.Dir(filepath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath, body, os.ModePerm); err != nil {
		return err
	}
	return nil
}
