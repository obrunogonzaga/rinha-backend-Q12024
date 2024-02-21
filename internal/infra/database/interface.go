package database

import (
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/entity"
)

type CustomerInterface interface {
	FindByID(id string) (*entity.Customer, error)
	Update(customer *entity.Customer) error
}

type TransactionInterface interface {
	CreateTransaction(transaction *entity.Transaction) error
}
