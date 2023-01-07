package main

import (
	"github.com/hardsky/gamma-beta/betta"
	"github.com/hardsky/gamma-beta/gamma"
	"log"
)

func main() {
	st, err := betta.NewPostgreStorage()
	if err != nil {
		log.Fatalf("database initialization: %s\n", err)
	}
	defer st.Close()

	gm := gamma.NewGamma()
	bt := betta.NewBetta(gm, st)
	bt.Run()
}
