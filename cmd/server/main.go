package main

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/infra/database"
	"github.com/obrunogonzaga/rinha-backend-Q12024/internal/infra/webserver/handlers"
	"net/http"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=admin password=123 dbname=rinha sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	customerDB := database.NewCustomer(db)

	customerHandler := handlers.NewCustomer(customerDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/clientes/{id}", customerHandler.GetCustomer)

	http.ListenAndServe(":8000", r)
}
