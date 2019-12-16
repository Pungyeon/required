package simple

import (
	"encoding/json"
	"fmt"
)

// User represents our user object, we wish to parse with our
// simple parsing method
type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Twitter   string `json:"twitter"`
	JobTitle  string `json:"job_title"`
}

// UserFromJSON will parse a json user and validate the required twitter property
func UserFromJSON(jsonUser []byte) (User, error) {
	var user User
	if err := json.Unmarshal(jsonUser, &user); err != nil {
		return User{}, err
	}
	if user.Twitter == "" {
		return User{}, fmt.Errorf("no value specified for required field user.twitter")
	}
	return user, nil
}
