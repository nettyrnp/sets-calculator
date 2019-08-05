package main

import (
	"reflect"
	"testing"
)

var (
	testdataDir = "testdata"
)

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("expecting no error")
	}
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Fatalf("expecting error")
	}
}

func assertEqual(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("expecting values to be equal but got: '%v' and '%v'", a, b)
	}
}

func TestGetFiles(t *testing.T) {
	files, err := getFiles(testdataDir)
	assertNoError(t, err)
	assertEqual(t, 3, len(files))
	assertEqual(t, file{"b.txt", []int{2, 3, 4}}, files[1])
}
