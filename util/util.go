package util

import (
	"github.com/nettyrnp/sets-calculator/log"
	"io/ioutil"
	"os"
)

// Die kills the failing program.
func Die(err error) {
	logger := log.GetLogger()
	if err == nil {
		return
	}
	logger.Fatal(err.Error())
	os.Exit(1)
}

func ReadFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
