package entity

type Transaction struct {
	ID              int    `json:"id"`
	CustomerID      int    `json:"customer_id"`
	Amount          int    `json:"amount"`
	TransactionType string `json:"transaction_type"`
	Description     string `json:"description"`
	CreatedAt       string `json:"created_at"`
}
