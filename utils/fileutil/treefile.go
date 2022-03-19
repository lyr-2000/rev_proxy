package fileutil

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func ListAllFilePathInDir(dir string) []string {
	//list
	var s []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		//fmt.Printf("%v, ", path)
		if info != nil && info.IsDir() == false {
			s = append(s, path)
		}
		return nil
	})
	if err != nil {
		log.Println("err =", err)
	}
	return s
}

func MkdirAll(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Printf("io error [%+v]\n", err)
	}
	return err
}
