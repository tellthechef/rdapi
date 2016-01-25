package main

import (
	rd "../"
	"fmt"
)

func main() {
	api := rd.New(0, "consumerKey", "consumerSecret", "secondSecret")
	if err := api.Authenticate(); err != nil {
		return
	}

	bookings, err := api.GetDiary("2016-01-21")
	if err != nil {
		panic(err)
	}

	for _, booking := range bookings {
		b, err := api.GetBooking(booking.ID)
		if err != nil {
			panic(err)
		}

		fmt.Println(b)
	}
}
