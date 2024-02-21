package handlers

import (
	"encoding/json"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/dto"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/infra/database"
	"log"
	"net/http"
)

type Transaction struct {
	TransactionDB database.TransactionInterface
}

func NewTransaction(transactionDB database.TransactionInterface) *Transaction {
	return &Transaction{
		TransactionDB: transactionDB,
	}
}

func (t *Transaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction dto.CreateTransactionInput
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Transaction: %+v", transaction)
	w.WriteHeader(http.StatusCreated)
}
