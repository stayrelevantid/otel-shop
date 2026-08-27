package model

// InventoryResponse is returned by GET /inventory/{id} (PRD §9).
type InventoryResponse struct {
	Stock int `json:"stock"`
}
