package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/dto"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/entity"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/infra/database"
	"log"
	"net/http"
	"strconv"
)

type Transaction struct {
	TransactionDB database.TransactionInterface
	CustomerDB    database.CustomerInterface
}

func NewTransaction(transactionDB database.TransactionInterface, customerDB database.CustomerInterface) *Transaction {
	return &Transaction{
		TransactionDB: transactionDB,
		CustomerDB:    customerDB,
	}
}

func (t *Transaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction dto.CreateTransactionInput
	id := chi.URLParam(r, "id")
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Transaction: %+v", transaction)

	customer, err := t.CustomerDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatal(err)
		return
	}
	log.Printf("Customer: %+v", customer)

	if transaction.TransactionType == "c" {
		customer.Balance += transaction.Amount
	} else {
		customer.Balance -= transaction.Amount
		limit := reverseSign(customer.Limit)
		if customer.Balance < limit {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	err = t.CustomerDB.Update(customer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	idInteger, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	trx, err := entity.NewTransaction(idInteger, transaction.Amount, transaction.TransactionType, transaction.Description)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = t.TransactionDB.CreateTransaction(trx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func reverseSign(num int) int {
	return num * -1
}
