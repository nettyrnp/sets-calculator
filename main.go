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

	task1 := expression{
		raw: `[ SUM a.txt b.txt c.txt ]`,
	}
	fmt.Printf(">> task: %v\n", task1)
	task2 := expression{
		raw: `[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]`,
	}
	fmt.Printf(">> task: %v\n", task2)
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

//type task struct {
//	raw      string
//	operator string
//	sets []set
//	output   string
//}
//
//func (t *task) parse() {
//	//parts:=[]string{}
//	// ...
//	// parsing ...
//	// ...
//	parts:=strings.Split(t.raw, " ")
//	t.operator=parts[0]
//	for _, part :=range parts[1:]{
//		t.sets=append(t.sets, set{
//			kind:       setKindExpr,
//			expression: expression{
//				raw: part,
//			},
//		})
//	}
//}
//
//
//func (t *task) do() {
//	t.output=""
//}

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
	sets     []set
	output   []int
}

func (e *expression) evaluate() {
	switch e.operator {
	case operatorKindSum:
		e.calcSum()
	}
}

func (e *expression) calcSum() {
	arr := []int{}
	for _, set := range e.sets {
		if set.kind == setKindFile {
			arr = append(arr, set.file...)
		} else if set.kind == setKindExpr {
			set.expression.evaluate()
		}
	}
	e.output = arr
}
