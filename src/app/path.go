package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var _workPath string

// fileExist 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func currentDirectory() string {
	if len(_workPath) > 0 {
		return _workPath
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	dir = strings.Replace(dir, "\\", "/", -1)
	checkDir := dir + "/resource"
	if !fileExist(checkDir) {
		dir, _ = os.Getwd()
		pos := strings.LastIndex(dir, "src")
		if pos == -1 {
			dir = ""
		} else {
			dir = dir[:pos]
		}
	}

	log.Println(dir)

	_workPath = dir
	return _workPath
}

func RootPath() string {
	return currentDirectory()
}

func Path(relativePath string) string {
	path := currentDirectory() + "/" + relativePath
	path = strings.Replace(path, "////", "/", -1)
	path = strings.Replace(path, "///", "/", -1)
	path = strings.Replace(path, "//", "/", -1)
	return path
}
