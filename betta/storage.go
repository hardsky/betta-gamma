package betta

import (
	"log"

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
	_, err := p.pool.Exec(context.Background(), "insert into votings(id, vote_id, option_id) values($1, $2, $3)", vote.VotingID, vote.VoteID, vote.OptionID)
	if err != nil {
		log.Printf("error on inserting vote, vote:%s, err:%s", vote, err)
	}
	return err
}

func (p *PostgresStorage) Close() {
	p.pool.Close()
}

func NewPostgreStorage() (Storage, error) {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
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
