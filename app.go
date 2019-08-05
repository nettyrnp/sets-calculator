package main

import (
	"fmt"
	"github.com/nettyrnp/sets-calculator/log"
	"github.com/nettyrnp/sets-calculator/util"
	"go.uber.org/zap"
	"regexp"
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
	task1.printResults()
	return nil
}
