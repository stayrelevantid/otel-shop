package model

// CheckoutRequest is the payload for POST /checkout (PRD §9).
type CheckoutRequest struct {
	ItemID string `json:"item_id"`
	Qty    int    `json:"qty"`
}

// CheckoutResponse is returned by POST /checkout on success.
type CheckoutResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

// ErrorResponse is the generic error body.
type ErrorResponse struct {
	Error string `json:"error"`
}

// PayRequest is sent by Order to the Payment service (PRD §8).
type PayRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

// InventoryResponse mirrors the Inventory service response (PRD §9).
type InventoryResponse struct {
	Stock int `json:"stock"`
}
