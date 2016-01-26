package models

import (
	"strconv"
	"strings"
	"time"
)

type ShortBooking struct {
	ID               int `json:"BookingId"`
	BookingReference string
	HasPayments      bool
	HasPromotions    bool
	CustomerFullName string
	CustomerID       int `json:"CustomerId"`
}

type Booking struct {
	ID               int `json:"BookingId"`
	ServiceID        int `json:"ServiceId"`
	RestaurantID     int `json:"RestaurantId"`
	AreaID           int `json:"AreaId"`
	BookingReference string
	BookingDateTime  string
	VisitDateTime    string
	Customer         Customer
	Duration         int
	Status           int
	Covers           int
	Comments         string
	ChannelName      string
	MenuName         string
	Type             int

	Promotions []BookingPromotion
	Payments   []BookingPayment
	Extras     []BookingExtra
}

func (booking *Booking) ParseBookingDate() *time.Time {
	if len(booking.VisitDateTime) == 0 {
		return nil
	}

	datetime := strings.Replace(strings.Replace(booking.VisitDateTime, "/Date(", "", 1), ")/", "", 1)
	dateInt, err := strconv.Atoi(datetime)
	if err != nil {
		return nil
	}

	date := time.Unix(int64(dateInt/1000), 0)
	return &date
}
