package betta

import "github.com/hardsky/gamma-beta/gamma"
import "github.com/hardsky/gamma-beta/models"
import "github.com/gin-gonic/gin"
import "net/http"
import "os/signal"
import "context"
import "syscall"
import "log"
import "time"

func NewBetta(clt gamma.GammaCollector) *Betta {
	return &Betta{clt}
}

type Betta struct {
	clt gamma.GammaCollector
}

func (p *Betta) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.POST("/voting", p.PostVote)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

// postVote process an vote from JSON received in the request body
func (p *Betta) PostVote(c *gin.Context) {
	var newVote models.Vote

	if err := c.BindJSON(&newVote); err != nil {
		return
	}

	c.IndentedJSON(http.StatusOK, models.NewResult("ok"))
}
