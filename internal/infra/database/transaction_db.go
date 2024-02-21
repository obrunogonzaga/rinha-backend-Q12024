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
	_, err := t.DB.Exec("INSERT INTO transactions (id, customer_id, amount, type) VALUES ($1, $2, $3, $4)", transaction.ID, transaction.CustomerID, transaction.Amount, transaction.TransactionType)
	if err != nil {
		return err
	}
	return nil
}
