package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return Store{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// func NewStore(db *sql.DB) Store {
// 	return Store{
// 		db:      db,
// 		Queries: New(db),
// 	}
// }
