package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db   *sql.DB
	repo *Repo
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Repo() *Repo {
	if s.repo != nil {
		return s.repo
	}
	s.repo = &Repo{
		store: s,
	}
	return s.repo
}
