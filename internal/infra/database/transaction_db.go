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

func (t Transaction) CreateTransaction(tx *sql.Tx, transaction *entity.Transaction) error {
	_, err := tx.Exec("INSERT INTO transactions (customer_id, amount, transaction_type, description) VALUES ($1, $2, $3, $4)", transaction.CustomerID, transaction.Amount, transaction.TransactionType, transaction.Description)
	if err != nil {
		return err
	}
	return nil
}

func (t Transaction) GetTransactionsByCustomerID(id string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	rows, err := t.DB.Query("SELECT id, customer_id, amount, transaction_type, description, timestamp FROM transactions WHERE customer_id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.ID, &transaction.CustomerID, &transaction.Amount, &transaction.TransactionType, &transaction.Description, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
