package main

import (
	"bytes"
	"github.com/nettyrnp/sets-calculator/util"
	"sort"
	"strings"
)

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
	m := util.ToMap(set1)
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

func (e *expression) calcDiff() {
	// store the first set
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

	// diff = "sum of all except first" DISJOIN "first set"
	resultSet := calcDiff0(sumSet, firstSet)

	sort.Ints(resultSet)
	e.output = resultSet
}

func calcDiff0(set1 []int, set2 []int) []int {
	arr := []int{}
	m := util.ToMap(set1)
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

	// extract the files/expressions
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
