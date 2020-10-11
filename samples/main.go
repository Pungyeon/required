package main

import (
	"github.com/Pungyeon/json-validation/pkg/json"
)

type Customer struct {
	ID      int32    `json:"id"`
	Name    string   `json:"name,required"`
	Address Address  `json:"address"`
	Tags    []string `json:"tags,required"`
}

type Address struct {
	StreetAddress1 string `json:"street_address_1,required"`
	StreetAddress2 string `json:"street_address_2"`
	Country        string `json:"country,required"`
	PostalCode     int    `json:"postal_code,required"`
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
	if err := json.Unmarshal(jsonBytes, &customer); err != nil {
		panic(err)
	}
}
