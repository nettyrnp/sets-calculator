package main

import (
	"testing"
)

func TestEvaluate(t *testing.T) {
	_, err := getFiles(testdataDir)
	assertNoError(t, err)

	testcases := []struct {
		in  string
		out []int
	}{
		{
			`[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]`,
			[]int{1, 3, 4},
		},
		{
			`[ INT c.txt b.txt ]`,
			[]int{3, 4},
		},
		{
			`[ INT b.txt c.txt ]`,
			[]int{3, 4},
		},
		{
			`[ INT a.txt b.txt c.txt ]`,
			[]int{3},
		},
		{
			`[ DIF b.txt a.txt ]`,
			[]int{4},
		},
		{
			`[ DIF a.txt b.txt ]`,
			[]int{1},
		},
		{
			`[ DIF a.txt b.txt c.txt ]`,
			[]int{1},
		},
		{
			`[ SUM [ DIF a.txt [ INT b.txt c.txt ] b.txt ] c.txt ]`,
			[]int{1, 3, 4, 5},
		},
		{
			`[ SUM [ DIF a.txt b.txt ] c.txt ]`,
			[]int{1, 3, 4, 5},
		},
	}

	for _, tc := range testcases {
		task := expression{}
		task.raw = tc.in
		task.evaluate()
		assertEqual(t, tc.out, task.output)
	}

}
