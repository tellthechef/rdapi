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
	CustomerSpend    int
	Duration         int
	Status           int
	Covers           int
	Comments         string
	ChannelID        int `json:"ChannelId"`
	ChannelName      string
	MenuId           int `json:"MenuId"`
	MenuName         string
	Type             int
	ArrivalStatus    int
	ConfirmedByPhone bool

	IsGuestIntendingToPayByApp bool
	IsLeaveTimeConfirmed       bool
	MealStatus                 int
	TurnTime                   int

	SpecialRequests []interface{}
	Promotions      []BookingPromotion
	Payments        []BookingPayment
	Extras          []BookingExtra
	Tables          []int
}

func (booking *Booking) ParseBookingDate() *time.Time {
	if len(booking.VisitDateTime) == 0 {
		return nil
	}

	datetime := strings.Replace(booking.VisitDateTime, "/Date(", "", 1)
	datetime = strings.Split(datetime, "+")[0]
	datetime = strings.Replace(datetime, ")/", "", 1)

	dateInt, err := strconv.ParseInt(datetime, 10, 64)
	if err != nil {
		return nil
	}

	date := time.Unix(int64(dateInt/1000), 0)
	return &date
}
