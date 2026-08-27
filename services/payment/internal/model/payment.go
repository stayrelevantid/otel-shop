package model

// PayRequest is the payload for POST /pay (PRD §8).
type PayRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

// PayResponse is returned by POST /pay.
type PayResponse struct {
	Status string `json:"status"`
}
