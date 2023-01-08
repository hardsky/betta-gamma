package models

import (
	"github.com/google/uuid"
)

// Vote represents data about user vote
type Vote struct {
	VoteID   uuid.UUID `json:"voteId"`
	VotingID uuid.UUID `json:"votingId"`
	OptionID uuid.UUID `json:"optionId"`
}

// StatisticVote represents staticstic for concrete voting
type StatisticVoting struct {
	VotingID uuid.UUID         `json:"votingId"`
	Results  []StatisticOption `json:"results"`
}

// StatisticOption represent statistic data for option in concrete voting
type StatisticOption struct {
	OptionID uuid.UUID `json:"optionId"`
	Count    int       `json:"count"`
}

// Result represent betta response
type Result struct {
	Result string `json:"result"`
}

func NewResult(msg string) *Result {
	return &Result{Result: msg}
}
