package main

import (
	"github.com/kavkaz/gol"
)

func main() {
	c := ParseArgs()
	gol.SetLevel(gol.DEBUG)
	gol.Debugf("config: %#v", *c)
}
