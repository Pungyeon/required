package validate

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

func (user User) Validate() error {
	return validateMany(
		validate(user.FirstName == "", "no value specified for required field user.first_name"),
		validate(user.LastName == "", "no value specified for required field user.last_name"),
		validate(user.Twitter == "", "no value specified for required field user.twitter"),
		validate(user.JobTitle == "", "no value specified for required field user.job_title"))
}

func validateMany(assertions ...error) error {
	for _, err := range assertions {
		if err != nil {
			return err
		}
	}
	return nil
}

func validate(assertion bool, format string, args ...interface{}) error {
	if assertion {
		return fmt.Errorf(format, args...)
	}
	return nil
}

// UserFromJSON will parse a json user and validate the required properties
func UserFromJSON(jsonUser []byte) (User, error) {
	var user User
	if err := validateMany(
		json.Unmarshal(jsonUser, &user),
		user.Validate(),
	); err != nil {
		return User{}, err
	}
	return user, nil
}

type Message struct {
	ID   int64 `json:"id"`
	To   User  `json:"to"`
	From User  `json:"from"`
}

// MessageFromJSON will parse a json message and validate the required properties
func MessageFromJSON(jsonMessage []byte) (Message, error) {
	var message Message
	if err := validateMany(
		json.Unmarshal(jsonMessage, &message),
		validate(message.ID == 0, "no value specified for required field message.id"),
		message.To.Validate(),
		message.From.Validate(),
	); err != nil {
		return Message{}, err
	}
	return message, nil
}
