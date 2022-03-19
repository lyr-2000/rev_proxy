package fileutil

import (
	"io/fs"
	"log"
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
