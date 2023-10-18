package dto

type TransactionsResponse struct {
	Title   string  `json:"title"`
	Message string  `json:"message"`
	Amount  float64 `json:"amount"`
	Type    int     `json:"type"`
}
