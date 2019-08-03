package main

import (
	"fmt"
	"github.com/nettyrnp/sets-calculator/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	setKindFile = iota + 1
	setKindExpr
)

const (
	operatorKindSum = "SUM"
	operatorKindInt = "INT"
	operatorKindDif = "DIF"
)

// set = file | expression
// expression = [ operator file_1 file_2 file_3 ..  file_N ]
// operator = SUM | INT | DIF
func main() {
	files, err := getFiles()
	util.Die(err)
	fmt.Printf(">> files: %v\n", files)

	task1 := `[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]`
	fmt.Printf(">> task: %v\n", task1)
}

func getFiles() ([]file, error) {
	files := []file{}
	rootDir, _ := os.Getwd()
	filesDir := filepath.Join(rootDir, "files")
	fmt.Printf(">> filesDir: %v\n", filesDir)

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
		files = append(files, file{f.Name(), arr})
	}
	return files, nil
}

func parseLines(body string) ([]int, error) {
	arr := []int{}
	for _, s := range strings.Split(body, "\n") { // todo: os-specific separator
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

type setKind int

type set struct {
	kind       setKind
	file       []int
	expression expression
}

type file struct {
	name string
	arr  []int
}

type expression struct {
	raw      string
	operator string
	operands []set
	output   []int
}

func (e *expression) evaluate() {
	switch e.operator {
	case operatorKindSum:
		e.output = calcSum(e.operands)
	}
}

func calcSum(operands []set) []int {
	out := []int{}
	for _, set := range operands {
		if set.kind == setKindFile {

		} else if set.kind == setKindExpr {

		}
	}
	return out
}
