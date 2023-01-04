package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// vote represents data about user vote
type vote struct {
	VoteID   uuid.UUID `json:"voteId"`
	VotingID uuid.UUID `json:"votingId"`
	OptionID uuid.UUID `json:"optionId"`
}

func main() {
	router := gin.Default()
	router.POST("/voting", postVote)

	router.Run("localhost:8080")
}

// postVote process an vote from JSON received in the request body
func postVote(c *gin.Context) {
	var newVote vote

	if err := c.BindJSON(&newVote); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, newVote)
}
