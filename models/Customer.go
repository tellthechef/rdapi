package models

type Customer struct {
	ID             int `json:"CustomerId"`
	CustomerTypeID int `json:"CustomerTypeId"`

	Title     string
	FirstName string
	Surname   string
	Email     string

	AccessCode        string
	Anniversary       string
	Birthday          string
	Comments          string
	Company           string
	Interests         []int
	RoyaltyCardHolder bool `json:"IsRoyaltyCardHolder"`
	VIP               bool `json:"IsVip"`

	MobileNumber string
	PhoneNumber  string
}
