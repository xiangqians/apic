// @author xiangqian
// @date 2025/08/16 02:09
package handler

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strings"
)

type Type string

const (
	TypeYaml Type = "yaml"
	TypeJson Type = "json"
	TypeJs   Type = "js"
)

type Dir struct {
	Id    string
	Name  string
	Files []File
}

type File struct {
	Id   string
	Name string
	Type Type
}

func readDirs() ([]Dir, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var dirs = make([]Dir, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			var name = entry.Name()
			files, err := readFiles(name)
			if err != nil {
				return nil, err
			}

			for _, file := range files {
				typ := file.Type
				if typ == TypeYaml || typ == TypeJson {
					dirs = append(dirs, Dir{Id: hash(name), Name: name, Files: files})
					break
				}
			}
		}
	}
	return dirs, nil
}

func readFiles(dir string) ([]File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files = make([]File, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			var name = entry.Name()
			var typ Type
			if strings.HasSuffix(name, ".yaml") {
				typ = TypeYaml
			} else if strings.HasSuffix(name, ".json") {
				typ = TypeJson
			} else if strings.HasSuffix(name, ".js") {
				typ = TypeJs
			}
			if typ != "" {
				files = append(files, File{Id: hash(name), Name: name, Type: typ})
			}
		}
	}
	return files, nil
}

func hash(data string) string {
	h := md5.Sum([]byte(data))
	return hex.EncodeToString(h[:])
}
