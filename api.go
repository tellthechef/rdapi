package rdapi

import (
	"encoding/json"
	"github.com/tellthechef/rdapi/models"
	"strconv"
	"time"
)

const DiaryDateFormat = "2006-01-02"

type RDConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	SecondSecret   string

	RestaurantID    int
	Endpoint        string
	ServiceEndpoint string

	firstAuth  authKeys
	secondAuth authKeys
}

func (api *RDConfig) GetDiary(date time.Time) ([]models.ShortBooking, error) {
	var bookings []models.ShortBooking

	client, req, _ := api.RestaurantRequest("GET", "/DiaryData?date="+date.Format(DiaryDateFormat), nil)
	res, _ := client.Do(req)

	err := json.NewDecoder(res.Body).Decode(&bookings)
	return bookings, err
}

func (api *RDConfig) GetBooking(id int) (models.Booking, error) {
	var booking models.Booking

	client, req, _ := api.RestaurantRequest("GET", "/Booking/"+strconv.Itoa(id), nil)
	res, _ := client.Do(req)

	err := json.NewDecoder(res.Body).Decode(&booking)
	return booking, err
}
