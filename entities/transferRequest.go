package entities

type TransferRequest struct {
	RecipientID int     `json:"recipient_id"`
	Amount      float64 `json:"amount"`
}
