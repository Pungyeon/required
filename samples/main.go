package main

import (
	"errors"
	"regexp"

	"github.com/Pungyeon/required/pkg/json"
)

type Customer struct {
	ID      int32    `json:"id"`
	Name    string   `json:"name,required"`
	Address Address  `json:"address"`
	Email   Email    `json:"email,required"`
	Tags    []string `json:"tags,required"`
}

type Address struct {
	StreetAddress1 string `json:"street_address_1,required"`
	StreetAddress2 string `json:"street_address_2"`
	Country        string `json:"country,required"`
	PostalCode     int    `json:"postal_code,required"`
}

type Email string

func (email Email) IsValueValid() error {
	matched, err := regexp.MatchString(`.+@.+\..+`, string(email))
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid email")
	}
	return nil
}

func main() {
	// Try deleting required fields from the JSON below,
	// and you should expect an error on line 39!
	jsonBytes := []byte(`{
		"name": "BigCustomer",
		"address": {
			"street_address_1": "Some General Name 24B",
			"country": "Dingeling",
			"postal_code": 91210
		},
		"email": "ding@dingeling.dk"
		"tags": ["big boi", "customer"]

	}`)

	var customer Customer
	if err := json.Unmarshal(jsonBytes, &customer); err != nil {
		panic(err)
	}
}
