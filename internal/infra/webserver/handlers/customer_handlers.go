package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/infra/database"
	"log"
	"net/http"
)

type Customer struct {
	CustomerDB database.CustomerInterface
}

func NewCustomer(customerDB database.CustomerInterface) *Customer {
	return &Customer{
		CustomerDB: customerDB,
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
	json.NewEncoder(w).Encode(customer)
}
