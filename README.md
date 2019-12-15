# Golang - Required JSON Fields 

## Introduction
So, recently at work, one our junior engineers asked me a question: "How do I create required fields for structures in Go, when parsing from JavaScript?". Now, I haven't done much work with API's in Go, so I'm actually not sure what the idiomatic solution for this is, however, delving into the topic turned out interesting. It's a perfect example of Go as an expressive language and how it allows you to approach problems from many different angles, despite it's sometimes restrictive nature. This article will describe some of these approached, by describing a few different techniques and methods for overcoming this issue without generics.

## The Simple Approach
So to begin with, let's discuss the most important part of solving this issues. We should never overcomplicate a solution. Always choose the solution, which is the most straight forward. In other words "Keep It Simple Stupid". So let's assume that our issue is that we have the following JSON object, we wish to parse:

```json
{
  "first_name": "Lasse",
  "last_name": "Jakobsen",
	"twitter": "ifndef_lmj",
  "job_title": "Lead Software Engineer"
}
```

In our first scenario, we wish to make sure that we have a value for the `twitter` property (for whatever reason). To ensure this, we will simply parse the given JSON and return either our parsed `User` object, or an error, if either the json couldn't be parsed, or if our required `twitter` value was not present.

```go
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

```

 If this is all that we need, this is the best solution. It's extremely simple and very straight forward. The code is particularly easy to understand and as we aren't overcomplicating the task. However, it's quite obvious how this is not quite as scalable, as none of this code is reusable. Our code may be simple, but as soon as we decide upon other required fields, for this or other structs, we will have to add more lines of code. If we decide that all of our fields are required, our code quickly becomes messy:

```go
// UserFromJSON will parse a json user and validate the required properties
func UserFromJSON(jsonUser []byte) (User, error) {
	var user User
	if err := json.Unmarshal(jsonUser, &user); err != nil {
		return User{}, err
	}
	if user.FirstName == "" {
		return User{}, fmt.Errorf("no value specified for required field user.first_name")
	}
	if user.LastName == "" {
		return User{}, fmt.Errorf("no value specified for required field user.last_name")
	}
	if user.Twitter == "" {
		return User{}, fmt.Errorf("no value specified for required field user.twitter")
	}
	if user.JobTitle == "" {
		return User{}, fmt.Errorf("no value specified for required field user.job_title")
	}
	return user, nil
}
```

## The Intermediate Approach
Fortunately for us, this is not too difficult to clean up. Again, sticking with the principle that we want to solve this issue as simply as possible, we can quickly create a new solution for handling our immediate task. If we look at our repeated code, we can quite quickly see the pattern. We check whether a `bool` value is true, and if it is, then we will return an error. Therefore, we can create a `validate` function for this:

```go
func validate(assertion bool, err error) error {
	if assertion {
		return err
	}
	return nil
}
```

This enables us to do a very small refactor, where we can substitute our assertions with the following:

```go
	...
	if err := validate(
		user.FirstName == "",
		fmt.Errorf("no value specified for required field user.first_name"),
	); err != nil {
		return User{}, err
	}
	...
```

Of course, doing this by itself, doesn't actually get us anywhere closer to cleaner code. We simply just wrapped our logical operation in a function and return an error, rather than the `bool`. In fact most would argue that this is actually worse, than what we had before. However, this is one of those cases where it's necessary to take a step back, before taking two steps forward. After this small refactor, we can see that all of our assertions are returning an `error` value. We can use this, to create another wrapping function, which will accept a _varadiac_ `error` parameter, which we will iterate, returning any eventual values, which aren't `nil`:

```go
func validateMany(assertions ...error) error {
	for _, err := range assertions {
		if err != nil {
			return err
		}
	}
	return nil
}

func validate(assertion bool, err error) error {
	if assertion {
		return err
	}
	return nil
}
```

In other words, our `validateMany` function is accepting `nil` or more `error` values, iterating through them and returning the `error` value, if it isn't `nil`. This allows us to significantly refactor our previous code:

```go
// UserFromJSON will parse a json user and validate the required properties
func UserFromJSON(jsonUser []byte) (User, error) {
	var user User
	if err := json.Unmarshal(jsonUser, &user); err != nil {
		return User{}, err
	}
	if err := validateMany(
		validate(user.FirstName == "", fmt.Errorf("no value specified for required field user.first_name")),
		validate(user.LastName == "", fmt.Errorf("no value specified for required field user.last_name")),
		validate(user.Twitter == "", fmt.Errorf("no value specified for required field user.twitter")),
		validate(user.JobTitle == "", fmt.Errorf("no value specified for required field user.job_title")),
	); err != nil {
		return User{}, err
	}
	return user, nil
}
```

Furthermore, we can change the function signature of our `validate` function, so that it accepts a `string` and `...interface{}` (like the `fmt.Errorf` function), so that we no longer need to repeat the `fmt.Errorf` on invoking `validate`:

```go
func validate(assertion bool, format string, args ...interface{}) error {
	if assertion {
		return fmt.Errorf(format, args...)
	}
	return nil
}
```

And, finally, our refactored code, looks as such:

```go
	...
	if err := validateMany(
		validate(user.FirstName == "", "no value specified for required field user.first_name"),
		validate(user.LastName == "", "no value specified for required field user.last_name"),
		validate(user.Twitter == "", "no value specified for required field user.twitter"),
		validate(user.JobTitle == "", "no value specified for required field user.job_title"),
	); err != nil {
		return User{}, err
	}
	...
```

Of course, we can also get rid of the repeated string `"no value specified for required field"` by wrapping our function once again:

```go
func required(assertion bool, field string) error {
  return validate(assertion, "no value specified for required field: %s", field)
}

// UserFromJSON will parse a json user and validate the required properties
func UserFromJSON(jsonUser []byte) (User, error) {
	var user User
	if err := json.Unmarshal(jsonUser, &user); err != nil {
		return User{}, err
	}
		...
	if err := validateMany(
		required(user.FirstName == "", "user.first_name"),
		required(user.LastName == "", "user.last_name"),
		required(user.Twitter == "", "user.twitter"),
		required(user.JobTitle == "", "user.job_title"),
	); err != nil {
		return User{}, err
	}
	return user, nil
}
```

A nice side effect of the last two refactors, is that we can use _any_ logical operator, with our `validate` function, to validate our required properties. We can now validate any type we wish, using this methodology. Furthermore, we aren't even coupled to our `validate` function, when using our `validateMany` function. Our `validateMany` function simply takes `...error` as input, so we can use any function which returns an error, in our validation. This opens up a lot of different options, as an example, we could extract our validation of the `User` into a method:

```go
// Validate all required fields of a given user
func (user User) Validate() error {
	return validateMany(
		required(user.FirstName == "", "user.first_name"),
		required(user.LastName == "", "user.last_name"),
		required(user.Twitter == "", "user.twitter"),
		required(user.JobTitle == "", "user.job_title"))	
}
```

And now, our `UserFromJSON` function can be refactored even further:

```go
// UserFromJSON will parse a json user and validate the required twitter property
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
```

Of course, this is not particularly useful in this case, but this now means that if a `User` is a property of another `struct`, it is now possible to easily validate the `User` property (as well as many others), by reusing the `validateMany` function. So, you can imagine future code looking something like this:

```go
type Message struct {
	ID int64
	  To User
	From User
  }
  
  func MessageFromJSON(jsonMessage []byte) (Message, error) {
	var message Message
	if err := validateMany(
		json.Unmarshal(jsonMessage, &message),
		message.To.Validate(),
		message.From.Validate(),
	); err != nil {
		return Message{}, err
	}
	return message, nil
  }
```

We now have a pretty solid solution, which is easy to implement for any of our structs with fields we wish to check as a requirement. However, this still isn't ideal. For each `struct` with a required field, we now have to implement a check on these fields. It would be much better if we could set a required tag on our structs. Unfortunately, despite this being proposed a few times, the Go programming language team have rejected this quite a few times now. Their argument has been both times, that they wish to keep the `json` package small and that it would be entirely possible to implement oneself.

## The Advanced Solution
I'm not going to lie. This won't be pretty, but I will try my best to show an example of how to implement this, in the form of a required struct type. We will start out, by showing an example of implementing a required string type. Let's start out with the easy part and create a wrapper struct for our `string` value:

```go
// String is a string type, which is required on JSON (un)marshal
type String struct {
	value string
}

// Value will return the inner string type
func (s *String) Value() string {
	return s.value
}
```

Next, we will need to define how to parse this from `json`, which will be a normal `string`, to our new `String` struct. To do this, we must implement the `json.Unmarshaler` interface on our struct, with the method `UnmarshalJSON`. This works, as the `json.Unmarshal` function, will eventually call this function, if the struct implements said interface. Our unmarshal function will look like this:

```go
var (
	// ErrStringEmpty represents an empty required string error
	ErrStringEmpty = errors.New("type of must.String not allowed to be empty")
	// ErrCannotUnmarshal represents an unmarshaling error
	ErrCannotUnmarshal = fmt.Errorf("json: cannot unmarshal into Go value of type must.String")
)


// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *String) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		if x == "" {
			return ErrStringEmpty
		}
		s.value = x
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
```

We aren't really doing anything crazy here. Essentially, we just use the standard `json.Unmarshal` into an `interface{}` and then checking the type of the unmarshalled json. If it turns out that it's a `string`, then we check if the value is equivalent of empty. If it is, then we return an error and if it's not, then we will assign the value to our `String.value` and return no error. Of course, if we find out it's not a `string` (which hopefully should never happen), we will return an error displaying our disgust.

 Of course, we also need to be able to marshal  `String` to JSON (and not only from). For this, we need to implement the `json.Marshaler` interface, by implementing the `MarshalJSON` method on our `String`. 

```go
// MarshalJSON is an implementation of the json.Marshaler interface
func (s String) MarshalJSON() ([]byte, error) {
	if s.Value() == "" {
		return []byte("null"), nil
	}
	return json.Marshal(s.value)
}
```

>  NOTE: Keep in mind, that we could also return an error if our `s.Value() == ""`.

Great! Easy. Let's just write some tests and then we are all good to go:

```go
type Person struct {
	Name String `json:"name"`
	Age  int64  `json:"age"`
}

func skipAssert(p Person) bool {
	return true
}

func TestStringValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(Person) bool
	}{
		{"valid strincg", `{"name":"Lasse"}`, nil, func(p Person) bool { return p.Name.Value() == "Lasse" }},
		{"empty string", `{"name":""}`, ErrStringEmpty, skipAssert},
		{"nil string", `{}`, ErrStringEmpty, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var person Person
			if err := json.Unmarshal(jsonb, &person); err != tf.err {
				t.Fatal(err)
			}

			if !tf.assert(person) {
				t.Fatalf("Assertion Failed: %+v", person)
			}
		})
	}
}
```

Running these tests however, gives us the following output:
```
--- FAIL: TestStringValidation (0.00s)
    --- FAIL: TestStringValidation/nil_string (0.00s)
        /Users/lassemartinjakobsen/projects/json-validation/required/string_test.go:34: <nil>
```

Wait? What the fuck? So it turns out, that life is just not that simple in the land of Go. The reason for this, is that when we pass the `"{}"` JSON string, the `json.Unmarshal` will ignore the fields which aren't present. This means that because the the `name` property doesn't appear in the JSON, the `UnmarshalJSON` function is never called, as it is completely skipped. Therefore, we never actually get a chance to check whether our required field actually contains anything, and therefore our tests fail miserably. 


