package main

import (
	"github.com/hardsky/gamma-beta/betta"
	"github.com/hardsky/gamma-beta/gamma"
)

func main() {
	gm := gamma.NewGamma()
	bt := betta.NewBetta(gm)
	bt.Run()
}
