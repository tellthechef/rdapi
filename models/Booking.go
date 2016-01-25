package models

type Booking struct {
	BookingID        int `json:"BookingId"`
	BookingReference string
	Customer         Customer
	Payments         []BookingPayment
}
