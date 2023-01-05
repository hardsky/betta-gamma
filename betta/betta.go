package betta

import "github.com/hardsky/gamma-beta/gamma"
import "github.com/hardsky/gamma-beta/models"
import "github.com/gin-gonic/gin"
import "net/http"

func NewBetta(clt gamma.GammaCollector) *Betta {
	return &Betta{clt}
}

type Betta struct {
	clt gamma.GammaCollector
}

func (p *Betta) Run() {
	router := gin.Default()
	router.POST("/voting", p.PostVote)

	router.Run("localhost:8080")

}

// postVote process an vote from JSON received in the request body
func (p *Betta) PostVote(c *gin.Context) {
	var newVote models.Vote

	if err := c.BindJSON(&newVote); err != nil {
		return
	}

	c.IndentedJSON(http.StatusOK, models.NewResult("ok"))
}
