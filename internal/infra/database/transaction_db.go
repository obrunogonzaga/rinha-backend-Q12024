package database

import (
	"database/sql"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/entity"
)

type Transaction struct {
	DB *sql.DB
}

func NewTransaction(db *sql.DB) *Transaction {
	return &Transaction{
		DB: db,
	}
}

func (t Transaction) CreateTransaction(transaction *entity.Transaction) error {
	_, err := t.DB.Exec("INSERT INTO transactions (customer_id, amount, transaction_type, description) VALUES ($1, $2, $3, $4)", transaction.CustomerID, transaction.Amount, transaction.TransactionType, transaction.Description)
	if err != nil {
		return err
	}
	return nil
}
