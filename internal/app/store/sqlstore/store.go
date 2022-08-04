package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db       *sql.DB
	userRepo *Repo
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Repo() *Repo {
	if s.userRepo != nil {
		return s.userRepo
	}
	s.userRepo = &Repo{
		store: s,
	}
	return s.userRepo
}
