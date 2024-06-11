package storage

import (
	"context"
	_ "embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	//go:embed query/tokens.sql
	tokensQuery string

	dsnLayout = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
)

type storage struct {
	db *sqlx.DB
}

func NewStorage() *storage {
	return &storage{}
}

func (s *storage) Open(opts StartOpts) error {
	dataSourceName := fmt.Sprintf(dsnLayout, opts.Host, opts.Port, opts.Username, opts.Password, opts.Database)
	db, err := sqlx.Open("pgx", dataSourceName)
	if err != nil {
		return err
	}

	if db.Ping() != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *storage) Close() error {
	return s.db.Close()
}

func (s *storage) Tokens(ctx context.Context) ([]string, error) {
	var dbTokens []token
	err := s.db.SelectContext(ctx, &dbTokens, tokensQuery)
	if err != nil {
		return nil, err
	}

	tokens := make([]string, 0, len(dbTokens))
	for _, t := range dbTokens {
		tokens = append(tokens, t.Token)
	}
	return tokens, nil
}
