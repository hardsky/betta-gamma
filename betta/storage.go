package betta

import (
	"github.com/sirupsen/logrus"

	"context"
	"os"

	"github.com/hardsky/gamma-beta/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	Store(vote models.Vote) error
	Close()
}

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func (p *PostgresStorage) Store(vote models.Vote) error {
	log := logrus.WithFields(logrus.Fields{
		"pkg":  "betta",
		"fnc":  "PostgresStorage.Store",
		"vote": vote,
	})
	_, err := p.pool.Exec(context.Background(), "insert into votings(id, vote_id, option_id) values($1, $2, $3)", vote.VotingID, vote.VoteID, vote.OptionID)
	if err != nil {
		log.WithField("err", err).Error("error on inserting vote")
	}
	return err
}

func (p *PostgresStorage) Close() {
	p.pool.Close()
}

func NewPostgreStorage() (Storage, error) {
	log := logrus.WithFields(logrus.Fields{
		"pkg": "betta",
		"fnc": "NewPostgreStorage",
	})
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.WithField("err", err).Error("error on database initialization")
		return nil, err
	}
	return &PostgresStorage{dbpool}, nil
}

/*
   create table votings(
     id UUID primary key,
     vote_id UUID not null,
     option_id UUID not null,
     created_at timestamptz not null default(now),
     unique(id, vote_id)
   );
*/
