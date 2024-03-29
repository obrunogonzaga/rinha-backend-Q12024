package database

import (
	"database/sql"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/entity"
	"log"
)

type Customer struct {
	DB *sql.DB
}

func NewCustomer(db *sql.DB) *Customer {
	return &Customer{
		DB: db,
	}
}

func (c Customer) FindByID(id string) (*entity.Customer, error) {
	log.Printf("Buscando cliente com id %s", id)
	var customer entity.Customer
	err := c.DB.QueryRow("SELECT id, limite, saldo FROM customers WHERE id = $1", id).Scan(&customer.ID, &customer.Limit, &customer.Balance)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (c Customer) Update(tx *sql.Tx, customer *entity.Customer) (*entity.Customer, error) {
	log.Printf("Atualizando cliente com id %s", customer.ID)
	row := tx.QueryRow("UPDATE customers SET saldo = $1 WHERE id = $2 RETURNING id, limite, saldo", customer.Balance, customer.ID)

	var updatedCustomer entity.Customer
	err := row.Scan(&updatedCustomer.ID, &updatedCustomer.Limit, &updatedCustomer.Balance)
	if err != nil {
		return nil, err
	}
	return &updatedCustomer, nil
}
