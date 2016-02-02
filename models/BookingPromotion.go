package models

type BookingPromotion struct {
	ID       int `json:"PromotionID"`
	Name     string
	Quantity int
}
