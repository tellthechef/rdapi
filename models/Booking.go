package models

type ShortBooking struct {
	ID               int `json:"BookingId"`
	BookingReference string
	HasPayments      bool
	HasPromotions    bool
}

type Booking struct {
	ID               int `json:"BookingId"`
	ServiceID        int `json:"ServiceId"`
	RestaurantID     int `json:"RestaurantId"`
	BookingReference string
	VisitDateTime    string `json:"VisitDateTime"`
	Customer         Customer
	CustomerSpend    int
	Duration         int
	Status           int

	Payments []BookingPayment
	Extras   []BookingExtra
}
