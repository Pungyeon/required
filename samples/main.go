package main

import (
	"fmt"

	"github.com/Pungyeon/json-validation/pkg/required"
)

type Customer struct {
	ID      int32                `json:"id"`
	Name    required.String      `json:"name"`
	Address Address              `json:"address"`
	Tags    required.StringSlice `json:"tags"`
}

type Address struct {
	StreetAddress1 required.String `json:"street_address_1"`
	StreetAddress2 string          `json:"street_address_2"`
	Country        required.String `json:"country"`
	PostalCode     required.Int    `json:"postal_code"`
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
		"tags": ["big boi", "customer"]

	}`)

	var customer Customer
	if err := required.Unmarshal(jsonBytes, &customer); err != nil {
		panic(err)
	}
	fmt.Println(customer)
}
