package db

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db DBTX
}

// NewStore creates a new Store
func NewStore(db DBTX) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
