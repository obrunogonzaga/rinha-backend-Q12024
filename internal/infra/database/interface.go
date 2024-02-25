package database

import (
	"database/sql"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/entity"
)

type CustomerInterface interface {
	FindByID(id string) (*entity.Customer, error)
	Update(tx *sql.Tx, customer *entity.Customer) (*entity.Customer, error)
}

type TransactionInterface interface {
	CreateTransaction(tx *sql.Tx, transaction *entity.Transaction) error
}
