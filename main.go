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
	fmt.Printf(">> files: %v\n\n", files)

	task1 := expression{
		//raw: `[ INT a.txt b.txt ]`,
		//raw: `[ INT b.txt c.txt ]`,
		//raw: `[ INT c.txt b.txt ]`,
		//raw: `[ INT a.txt b.txt c.txt ]`,
		//raw: `[ DIF a.txt b.txt ]`,
		//raw: `[ DIF b.txt a.txt ]`,
		//raw: `[ DIF a.txt b.txt c.txt ]`,
		//raw: `[ SUM a.txt b.txt c.txt ]`,
		//raw: `[ SUM a.txt  [ SUM a.txt b.txt ] b.txt ]`,
		//raw: `[ SUM c.txt  [ INT a.txt b.txt c.txt ] ]`,
		//raw: `[ SUM [ DIF a.txt b.txt ] [ INT b.txt c.txt ] ]`,
		raw: `[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]`, // 1 + 3,4 = 1,3,4
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
	raw         string
	operator    string
	sets        []set
	output      []int
	isEvaluated bool
}

func (e *expression) evaluate() {
	if e.isEvaluated {
		return
	}
	if e.operator == "" && e.sets == nil {
		e.parse()
	}
	switch e.operator {
	case operatorKindSum:
		e.calcSum()
	case operatorKindInt:
		e.calcIntersection()
	case operatorKindDif:
		e.calcDiff()
	}
	e.isEvaluated = true
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
			if !set.expression.isEvaluated {
				set.expression.evaluate()
			}
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
	// store the first set to compare the others against it
	var resultSet []int
	m := map[int]struct{}{}
	for _, set := range e.sets {
		if set.kind == setKindFile {
			if resultSet == nil {
				resultSet = set.file
			}
		} else if set.kind == setKindExpr {
			if !set.expression.isEvaluated {
				set.expression.evaluate()
			}
			if resultSet == nil {
				resultSet = set.expression.output
			}
		}
	}
	// calculate intersection pair-wise
	for i, set := range e.sets {
		if i > 0 {
			if set.kind == setKindFile {
				resultSet = calcIntersection0(resultSet, set.file)
			} else if set.kind == setKindExpr {
				if !set.expression.isEvaluated {
					set.expression.evaluate()
				}
				resultSet = calcIntersection0(resultSet, set.expression.output)
			}
		}
	}

	for k, _ := range m {
		resultSet = append(resultSet, k)
	}
	sort.Ints(resultSet)
	e.output = resultSet
}

func calcIntersection0(set1 []int, set2 []int) []int {
	arr := []int{}
	m := toMap(set1)
	resMap := map[int]struct{}{}
	for _, k := range set2 {
		if _, ok := m[k]; ok {
			resMap[k] = struct{}{}
		}
	}
	for k, _ := range resMap {
		arr = append(arr, k)
	}
	return arr
}

// diff = sum of all except first DISJOIN first
func (e *expression) calcDiff() {
	// store the first set to compare the others against it
	var firstSet []int
	set := e.sets[0]
	if set.kind == setKindFile {
		firstSet = set.file
	} else if set.kind == setKindExpr {
		if !set.expression.isEvaluated {
			set.expression.evaluate()
		}
		firstSet = set.expression.output
	}

	// store the sum of all sets except first one
	newExpr := e
	newExpr.operator = operatorKindSum
	b := append(e.sets[:0:0], e.sets...)
	newExpr.sets = b[1:]
	newExpr.calcSum()
	sumSet := newExpr.output

	resultSet := calcDiff0(sumSet, firstSet)

	sort.Ints(resultSet)
	e.output = resultSet
}

func calcDiff0(set1 []int, set2 []int) []int {
	arr := []int{}
	m := toMap(set1)
	resMap := map[int]struct{}{}
	for _, k := range set2 {
		if _, ok := m[k]; !ok {
			resMap[k] = struct{}{}
		}
	}
	for k, _ := range resMap {
		arr = append(arr, k)
	}
	return arr
}

func (e *expression) parse() {
	// extract the operator
	raw := e.raw
	if !re.MatchString(raw) {
		logger.Fatalf("%v is not a valid expression", raw)
	}
	gr := re.FindStringSubmatch(raw)
	head := gr[2]
	tail := gr[3]
	e.operator = head

	// extract the expressions
	cur, prev := 0, 0
	left, right := "[", "]"
	buf := bytes.Buffer{}
	arr := []set{}
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
				arr = append(arr, set{
					kind: setKindFile,
					file: filesMap[filename].arr,
				})
				buf.Reset()
			}
		}
		if prev == 1 && cur == 0 {
			expr := strings.TrimSpace(buf.String())
			arr = append(arr, set{
				kind: setKindExpr,
				expression: expression{
					raw: expr,
				},
			})
			buf.Reset()
		}
		prev = cur
	}
	e.sets = arr
}

func toMap(set []int) map[int]struct{} {
	m := map[int]struct{}{}
	for _, k := range set {
		m[k] = struct{}{}
	}
	return m
}
