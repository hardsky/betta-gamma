package gamma

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hardsky/gamma-beta/models"
	"github.com/sirupsen/logrus"
)

func NewGamma() GammaCollector {
	return &MemmoryGammaCollector{make(map[uuid.UUID]map[uuid.UUID]int), &HTTPGammaSender{}}
}

// GammaCollector collects votes for Gamma
type GammaCollector interface {
	Vote(vote models.Vote) error
}

type GammaSender interface {
	Send(statistic models.StatisticVoting) error
}

type MemmoryGammaCollector struct {
	votes  map[uuid.UUID]map[uuid.UUID]int
	sender GammaSender
}

func (p *MemmoryGammaCollector) Vote(vote models.Vote) error {
	if p.willVotingChanged(vote) {
		var voting map[uuid.UUID]int
		voting, ok := p.votes[vote.VotingID]
		if !ok {
			voting = make(map[uuid.UUID]int)
			p.votes[vote.VotingID] = voting
			voting[vote.OptionID] = 1
		}

		stats := models.StatisticVoting{VotingID: vote.VotingID}
		for k, v := range voting {
			stats.Results = append(stats.Results, models.StatisticOption{k, v})
		}

		return p.sender.Send(stats)
	}
	return nil
}

func (p *MemmoryGammaCollector) willVotingChanged(vote models.Vote) bool {
	var voting map[uuid.UUID]int
	if _, ok := p.votes[vote.VotingID]; !ok {
		return true
	}
	voting = p.votes[vote.VotingID]

	if _, ok := voting[vote.OptionID]; !ok {
		return true
	}

	var total int = 0
	for _, v := range voting {
		total += v
	}

	percents := make(map[uuid.UUID]int)
	for k, v := range voting {
		percents[k] = v / total
	}

	percents2 := make(map[uuid.UUID]int)
	opt := voting[vote.OptionID]
	voting[vote.OptionID] = (opt + 1)

	for k, v := range voting {
		percents2[k] = v / total
	}

	for k, _ := range voting {
		if percents[k] != percents2[k] {
			return true
		}
	}

	return false
}

type HTTPGammaSender struct {
}

func (p *HTTPGammaSender) Send(statistic models.StatisticVoting) error {
	log := logrus.WithFields(logrus.Fields{
		"pkg":  "gamma",
		"fnc":  "HTTPGammaSender.Send",
		"stats": statistic,
	})
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	reqBytes, err := json.Marshal(statistic)
	req, err := http.NewRequest("POST", "http://service-gamma/voting-stats/", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.WithField("err", err).Error("error on posting data to gamma")
		return err
	}
	
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.WithField("err", err).Error("error on posting data to gamma")
		return err
	}
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithField("err", err).Error("error on gamma response")
		return err
	}

	return err
}
