package storage

import (
	"sync"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db    *sqlx.DB
	mutex *sync.Mutex
}

func New(db *sqlx.DB) *Storage {
	return &Storage{
		db:    db,
		mutex: new(sync.Mutex),
	}
}
