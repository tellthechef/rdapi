package models

type BookingPayment struct {
	Amount        float64
	Processed     bool `json:"IsProcessed"`
	PaymentMethod string
	ProcessedOn   string
}
