package dto

type CreateTransactionInput struct {
	Amount          int    `json:"valor"`
	TransactionType string `json:"tipo"`
	Description     string `json:"descricao"`
}
