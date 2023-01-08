package betta

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hardsky/gamma-beta/gamma"
	"github.com/hardsky/gamma-beta/models"
	"github.com/sirupsen/logrus"
)

func NewBetta(clt gamma.GammaCollector, store Storage) *Betta {
	return &Betta{clt, store}
}

type Betta struct {
	clt gamma.GammaCollector
	st  Storage
}

func (p *Betta) Run() {
	log := logrus.WithFields(logrus.Fields{
		"pkg": "betta",
		"fnc": "Betta.Run",
	})

	log.Info("server is starting...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.POST("/voting", p.PostVote)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log = log.WithFields(logrus.Fields{"addr": srv.Addr})

	log.Info("server is working now")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithFields(logrus.Fields{"err": err}).Fatal("server stop listen")
		}
	}()

	log.Info("stop server pass Ctrl+C")
	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithFields(logrus.Fields{"err": err}).Fatal("Server forced to shutdown")
	}

	log.Info("server exiting")
}

// postVote process an vote from JSON received in the request body
func (p *Betta) PostVote(c *gin.Context) {
	var newVote models.Vote

	if err := c.BindJSON(&newVote); err != nil {
		return
	}

	go func(vote models.Vote) {
		p.clt.Vote(vote)
		p.st.Store(vote)
	}(newVote)

	c.IndentedJSON(http.StatusOK, models.NewResult("ok"))
}
