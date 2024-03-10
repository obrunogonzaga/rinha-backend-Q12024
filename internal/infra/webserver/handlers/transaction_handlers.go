package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/dto"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/entity"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/infra/database"
	"log"
	"net/http"
	"strconv"
	"sync"
)

const (
	CreditTransaction = "c"
	DebitTransaction  = "d"
)

type Transaction struct {
	DB            *sql.DB
	TransactionDB database.TransactionInterface
	CustomerDB    database.CustomerInterface
}

type Response struct {
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}

func NewTransaction(db *sql.DB, transactionDB database.TransactionInterface, customerDB database.CustomerInterface) *Transaction {
	return &Transaction{
		DB:            db,
		TransactionDB: transactionDB,
		CustomerDB:    customerDB,
	}
}

func (t *Transaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction dto.CreateTransactionInput
	id := chi.URLParam(r, "id")

	idInteger, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	customer, err := t.CustomerDB.FindByID(id)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	err = t.updateCustomerBalance(customer, transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	tx, err := t.DB.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	err = t.updateCustomerAndCreateTransaction(tx, idInteger, customer, transaction)
	if err != nil {
		http.Error(w, "Failed to update or create transaction", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	response := Response{
		Limit:   customer.Limit,
		Balance: customer.Balance,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (t *Transaction) GetTransactionsByCustomerId(id string) ([]entity.Transaction, error) {
	if id == "" {
		log.Fatal("Invalid ID")
		return nil, errors.New("Invalid ID")
	}

	transactions, err := t.TransactionDB.GetTransactionsByCustomerID(id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return transactions, nil
}

func (t *Transaction) updateCustomerAndCreateTransaction(tx *sql.Tx, customerId int, customer *entity.Customer, transaction dto.CreateTransactionInput) error {
	var wg sync.WaitGroup
	var customerErr, transactionErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		customer, customerErr = t.CustomerDB.Update(tx, customer)
	}()

	go func() {
		defer wg.Done()
		trx, err := entity.NewTransaction(customerId, transaction.Amount, transaction.TransactionType, transaction.Description)
		if err != nil {
			transactionErr = err
			return
		}
		transactionErr = t.TransactionDB.CreateTransaction(tx, trx)
	}()

	wg.Wait()

	if customerErr != nil {
		return customerErr
	}

	if transactionErr != nil {
		return transactionErr
	}

	return nil
}

func (t *Transaction) updateCustomerBalance(customer *entity.Customer, transaction dto.CreateTransactionInput) error {
	if transaction.TransactionType == CreditTransaction {
		customer.Balance += transaction.Amount
	} else if transaction.TransactionType == DebitTransaction {
		customer.Balance -= transaction.Amount
		limit := reverseSign(customer.Limit)
		if customer.Balance < limit {
			return errors.New("Insufficient balance")
		}
	} else {
		return errors.New("Invalid transaction type")
	}
	return nil
}

func reverseSign(num int) int {
	return num * -1
}
