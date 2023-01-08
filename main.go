package main

import (
	"github.com/sirupsen/logrus"

	"github.com/hardsky/gamma-beta/betta"
	"github.com/hardsky/gamma-beta/gamma"
)

func main() {

	log := logrus.WithFields(logrus.Fields{
		"pkg": "main",
		"fnc": "main",
	})

	log.Info("initialize storage")
	st, err := betta.NewPostgreStorage()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("storage initialization failed")
	}
	defer st.Close()

	log.Info("initialize gamma")
	gm := gamma.NewGamma()

	log.Info("initialize betta")
	bt := betta.NewBetta(gm, st)

	log.Info("start betta")
	bt.Run()

	log.Info("betta have completed work")
}

