package main

import (
	"fmt"
	"github.com/nettyrnp/sets-calculator/log"
	"github.com/nettyrnp/sets-calculator/util"
	"go.uber.org/zap"
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
	logger   = log.GetLogger()
	re       = regexp.MustCompile(fmt.Sprintf(`\[ *((%v|%v|%v) +(.+)) *\]`, operatorKindSum, operatorKindInt, operatorKindDif))
	filesMap = map[string]file{}
)

type App struct {
	config Config
	logger *zap.SugaredLogger
}

func NewApp() (*App, error) {
	logger := log.GetLogger()
	defer logger.Sync() // flushes buffer, if any

	conf, err := GetConfig()
	util.Die(err)

	a := &App{
		config: *conf,
		logger: logger,
	}
	return a, nil
}

func (a *App) Run() error {
	_, err := getFiles(a.config.Folder)
	util.Die(err)
	task1 := expression{
		raw: a.config.Task,
	}
	task1.evaluate()
	fmt.Printf(">> Task: %v\n\n", task1)
	fmt.Printf(">> Task.output: %v\n\n", task1.output)

	return nil
}

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
//raw: `[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]`, // 1 + 3,4 = 1,3,4
//raw: `[ SUM [ DIF a.txt b.txt c.txt ] a.txt b.txt [ INT b.txt c.txt ] c.txt ]`,
//raw: `[ SUM [ DIF a.txt [ INT b.txt c.txt ] b.txt ] c.txt ]`,
//raw: `[ SUM [ DIF a.txt b.txt ] c.txt ]`,
//raw: `[ SUM c.txt [ DIF a.txt b.txt ] ]`,
