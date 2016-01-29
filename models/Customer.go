package models

import "strings"

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

func (customer *Customer) GetFullName() string {
	return strings.Join([]string{customer.Title, customer.FirstName, customer.Surname}, " ")
}
