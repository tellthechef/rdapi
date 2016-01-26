package main

import (
	rd "../"
	"fmt"
	"time"
)

func main() {
	api := rd.New(0, "consumerKey", "consumerSecret", "secondSecret")
	if err := api.Authenticate(); err != nil {
		return
	}

	bookings, err := api.GetDiary(time.Date(2016, 1, 21, 0, 0, 0, 0, time.Local))
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
