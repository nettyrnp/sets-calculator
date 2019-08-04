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
	"sort"
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
	logger   = log.GetLogger()
	re       = regexp.MustCompile(fmt.Sprintf(`\[ *((%v|%v|%v) +(.+)) *\]`, operatorKindSum, operatorKindInt, operatorKindDif))
	filesMap = map[string]file{}
)

// set = file | expression
// expression = [ operator file_1 file_2 file_3 ..  file_N ]
// operator = SUM | INT | DIF
func main() {
	files, err := getFiles()
	util.Die(err)
	fmt.Printf(">> files: %v\n", files)

	task1 := expression{
		//raw: `[ SUM a.txt b.txt c.txt ]`,
		raw: `[ SUM a.txt  [ SUM a.txt b.txt ] b.txt ]`,
		//	//raw: `[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]`,
		//	raw: `[ SUM [ DIF a.txt b.txt c.txt ] a.txt b.txt [ INT b.txt c.txt ] c.txt ]`,
		//	//raw: `[ SUM [ DIF a.txt [ INT b.txt c.txt ] b.txt ] c.txt ]`,
		//	//raw: `[ SUM [ DIF a.txt b.txt ] c.txt ]`,
		//	//raw: `[ SUM c.txt [ DIF a.txt b.txt ] ]`,
	}
	task1.evaluate()
	fmt.Printf(">> task: %v\n\n", task1)
	fmt.Printf(">> task.output: %v\n\n", task1.output)
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
		newFile := file{f.Name(), arr}
		files = append(files, newFile)
		filesMap[f.Name()] = newFile
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
	if e.operator == "" && e.sets == nil {
		e.parse()
	}
	switch e.operator {
	case operatorKindSum:
		e.calcSum()
	}
}

func (e *expression) calcSum() {
	arr := []int{}
	m := map[int]struct{}{}
	for _, set := range e.sets {
		if set.kind == setKindFile {
			for _, i := range set.file {
				m[i] = struct{}{}
			}
		} else if set.kind == setKindExpr {
			set.expression.evaluate()
			for _, i := range set.expression.output {
				m[i] = struct{}{}
			}
		}
	}
	for k, _ := range m {
		arr = append(arr, k)
	}
	sort.Ints(arr)
	e.output = arr
}

func (e *expression) calcIntersection() {
	arr := []int{}
	m := map[int]struct{}{}
	for _, set := range e.sets {
		if set.kind == setKindFile {
			for _, i := range set.file {
				m[i] = struct{}{}
			}
		} else if set.kind == setKindExpr {
			set.expression.evaluate()
			for _, i := range set.expression.output {
				m[i] = struct{}{}
			}
		}
	}
	for k, _ := range m {
		arr = append(arr, k)
	}
	sort.Ints(arr)
	e.output = arr
}

//func sum(sets ...int) []int {
//	arr := []int{}
//	m:=map[int]struct{}{}
//	for _,st:=range sets{
//		m[]
//	}
//	return arr
//}

func (e *expression) parse() {
	// parse the operator
	raw := e.raw
	if !re.MatchString(raw) {
		logger.Fatalf("%v is not a valid expression", raw)
	}
	gr := re.FindStringSubmatch(raw)
	head := gr[2]
	tail := gr[3]
	fmt.Printf("tail: '%v'\n", tail)
	e.operator = head

	// parse the tail
	cur, prev := 0, 0
	left, right := "[", "]"
	buf := bytes.Buffer{}
	//arr := []string{}
	arr2 := []set{}
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
				filename := strings.TrimSpace(buf.String())
				//arr = append(arr, filename)
				arr2 = append(arr2, set{
					kind: setKindFile,
					file: filesMap[filename].arr,
				})
				buf.Reset()
			}
		}
		if prev == 1 && cur == 0 {
			expr := strings.TrimSpace(buf.String())
			//arr = append(arr, expr)
			arr2 = append(arr2, set{
				kind: setKindExpr,
				expression: expression{
					raw: expr,
				},
			})
			buf.Reset()
		}
		prev = cur
	}
	//fmt.Printf("arr: %v\n", arr)
	fmt.Printf("arr2: %v\n", arr2)
	e.sets = arr2

	//// parse sub-expressions recursively
	//for _,st :=range e.sets{
	//
	//}

	//switch head {
	//case operatorKindSum:
	//	e.operator=head
	//}

	//for _, part :=range parts[1:]{
	//	e.sets=append(e.sets, set{
	//		kind:       setKindExpr,
	//		expression: expression{
	//			raw: part,
	//		},
	//	})
	//}
}
