package platform

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func NewRiver(dbpool *pgxpool.Pool, opts *river.Config) (*river.Client[pgx.Tx], error) {
	client, err := river.NewClient(riverpgxv5.New(dbpool), opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}
