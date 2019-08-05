package main

import (
	"github.com/nettyrnp/sets-calculator/util"
)

func main() {
	app, err := NewApp()
	util.Die(err)
	util.Die(app.Run())
}
