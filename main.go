package main

import (
	"bytes"
	"fmt"
	"github.com/nettyrnp/sets-calculator/log"
	"github.com/nettyrnp/sets-calculator/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

var (
	logger = log.GetLogger()
	re     = regexp.MustCompile(fmt.Sprintf(`\[ *((%v|%v|%v) +(.+)) *\]`, operatorKindSum, operatorKindInt, operatorKindDif))
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
	task1.evaluate()
	fmt.Printf(">> task: %v\n\n", task1)

	task2 := expression{
		//raw: `[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]`,
		raw: `[ SUM [ DIF a.txt b.txt c.txt ] a.txt b.txt [ INT b.txt c.txt ] c.txt ]`,
		//raw: `[ SUM [ DIF a.txt b.txt ] c.txt ]`,
		//raw: `[ SUM c.txt [ DIF a.txt b.txt ] ]`,
	}
	task2.evaluate()
	fmt.Printf(">> task: %v\n\n", task2)
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
	sets     []set
	output   []int
}

func (e *expression) evaluate() {
	if e.operator == "" {
		e.parse()
	}
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

func (e *expression) parse() {
	raw := e.raw
	if !re.MatchString(raw) {
		logger.Fatalf("%v is not a valid expression", raw)
	}
	gr := re.FindStringSubmatch(raw)
	head := gr[2]
	tail := gr[3]
	fmt.Printf("tail: '%v'\n", tail)
	e.operator = head

	// proceed
	cur, prev := 0, 0
	left, right := "[", "]"
	buf := bytes.Buffer{}
	arr := []string{}
	for i := 0; i < len(tail); i++ {
		s := tail[i : i+1]
		if s == left {
			cur++
		}
		if s == right {
			cur--
		}
		buf.WriteString(s)
		if prev == 0 && cur == 0 {
			if s == " " {
				arr = append(arr, buf.String())
				buf.Reset()
			}
		}
		if prev == 1 && cur == 0 {
			arr = append(arr, buf.String())
			buf.Reset()
		}
		prev = cur
	}
	fmt.Printf("arr: %v\n", arr)

	//switch head {
	//case operatorKindSum:
	//	e.operator=head
	//}

	////parts:=[]string{}
	//// todo ...
	//parts:=strings.Split(e.raw, " ")
	//e.operator=parts[0]
	//for _, part :=range parts[1:]{
	//	e.sets=append(e.sets, set{
	//		kind:       setKindExpr,
	//		expression: expression{
	//			raw: part,
	//		},
	//	})
	//}
}
