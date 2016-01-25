package rdapi

import (
	"encoding/json"
	"github.com/tellthechef/rdapi/models"
)

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

func (api *RDConfig) GetDiary(date string) ([]models.Booking, error) {
	var bookings []models.Booking

	client, req, _ := api.RestaurantRequest("GET", "/DiaryData?date="+date, nil)
	res, _ := client.Do(req)

	err := json.NewDecoder(res.Body).Decode(&bookings)
	return bookings, err
}
