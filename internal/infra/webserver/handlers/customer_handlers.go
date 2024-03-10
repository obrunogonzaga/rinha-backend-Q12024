package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/infra/database"
	"log"
	"net/http"
	"time"
)

type Customer struct {
	DB            *sql.DB
	CustomerDB    database.CustomerInterface
	TransactionDB database.TransactionInterface
}

type Saldo struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}

type UltimasTransacoes []struct {
	Valor       int    `json:"valor"`
	Tipo        string `json:"tipo"`
	Descricao   string `json:"descricao"`
	RealizadaEm string `json:"realizada_em"`
}

type CustomerOutput struct {
	Saldo             Saldo             `json:"saldo"`
	UltimasTransacoes UltimasTransacoes `json:"ultimas_transacoes"`
}

func NewCustomer(db *sql.DB, customerDB database.CustomerInterface, transactionDB database.TransactionInterface) *Customer {
	return &Customer{
		DB:            db,
		CustomerDB:    customerDB,
		TransactionDB: transactionDB,
	}
}

func (c *Customer) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	customer, err := c.CustomerDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func (c *Customer) GetStatement(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	customer, err := c.CustomerDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatal(err)
		return
	}

	balanceOutput := Saldo{
		Total:       customer.Limit,
		Limite:      customer.Limit,
		DataExtrato: time.Now().UTC(),
	}

	// Get last transactions in database for account
	transactions, err := c.TransactionDB.GetTransactionsByCustomerID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	var lastTransactions UltimasTransacoes
	for _, transaction := range transactions {
		lastTransactions = append(lastTransactions, struct {
			Valor       int    `json:"valor"`
			Tipo        string `json:"tipo"`
			Descricao   string `json:"descricao"`
			RealizadaEm string `json:"realizada_em"`
		}{
			Valor:       transaction.Amount,
			Tipo:        transaction.TransactionType,
			Descricao:   transaction.Description,
			RealizadaEm: transaction.CreatedAt,
		})
	}

	customerOutput := CustomerOutput{
		Saldo:             balanceOutput,
		UltimasTransacoes: lastTransactions,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customerOutput)
}
