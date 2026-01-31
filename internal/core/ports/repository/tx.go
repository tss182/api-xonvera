package portRepository

// Transaction represents a database transaction
type Transaction interface {
	Commit() error
	Rollback() error
}

// TxRepository manages database transactions
type TxRepository interface {
	Begin() (Transaction, error)
}
