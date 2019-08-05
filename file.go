package main

import (
	"github.com/nettyrnp/sets-calculator/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getFiles(dir string) ([]file, error) {
	files := []file{}
	rootDir, _ := os.Getwd()
	filesDir := filepath.Join(rootDir, dir)

	fileInfo, err := ioutil.ReadDir(filesDir)
	if err != nil {
		return files, err
	}

	for _, f := range fileInfo {
		if f.IsDir() {
			continue
		}
		body, err := util.ReadFile(filesDir + "/" + f.Name())
		if err != nil {
			return files, err
		}
		arr, err := parseLines(body)
		newFile := file{f.Name(), arr}
		files = append(files, newFile)
		filesMap[f.Name()] = newFile
	}
	return files, nil
}

func parseLines(body string) ([]int, error) {
	arr := []int{}
	for _, s := range strings.Split(body, "\n") { // todo: os-independent separator
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return arr, err
		}
		arr = append(arr, v)
	}
	return arr, nil
}
