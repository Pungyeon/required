# Golang - Required JSON Fields 

## Introduction
So, recently at work, one our junior engineers asked me a question: "How do I create required fields for structures in Go, when parsing from JavaScript?". Now, I haven't done too much work with API's in Go, so I'm actually not sure what the idiomatic solution for this is, however, delving into the topic turned out interesting. It's a perfect example of Go as an expressive language and how it allows you to approach problems from many different angles, despite it's sometimes restrictive nature. This article will describe some of these approached, by describing a few different techniques and methods for overcoming this issue without generics.

## Where to find the code
You are probably already on github, but in case you are not, you can find all the final code of each section here: https://github.com/Pungyeon/required/article

## The Simple Approach
So to begin with, let's discuss the most important part of solving this issue. We should never overcomplicate a solution. Always choose the solution, which is the most straight forward. In other words "Keep It Simple Stupid". So let's start by solving the simplest issue, in the simplest manner. Consider the following JSON object:

```json
{
  "first_name": "Lasse",
  "last_name": "Jakobsen",
	"twitter": "ifndef_lmj",
  "job_title": "Lead Software Engineer"
}
```

In our first scenario, we wish to make sure that we have a value for the `twitter` property (for whatever reason). To ensure this, we will simply parse the given JSON and return either our parsed `User` object, or an error if either the json couldn't be parsed, or if our required `twitter` value was not present.

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

 If this is all that we need, this is the best solution. It's extremely simple and very straight forward. The code is particularly easy to understand and as we aren't overcomplicating the task. 
 
 However, it's quite apparent how this solution is not particularly scalable, as none of this code is reusable. Our code may be simple, but as soon as we decide upon adding other required fields, for this or other structs, we will have to add more lines of code. If we decide that all of our fields are required, our code quickly becomes messy:

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
Fortunately for us, this is not too difficult to clean up. Again, sticking with the principle that we wish to solve this issue as simply as possible, we can quickly create a new solution for handling our immediate task. If we look at our repeated code, we can quite quickly see the pattern. We check whether a `bool` value is true, and if it is, then we will return an error. Therefore, we can create a `validate` function for this:

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

Of course, doing this by itself, doesn't actually get us anywhere closer to cleaner code. We simply just wrapped our logical operation in a function and return an error, rather than the `bool`. In fact most would argue that this is actually worse, than what we had before. However, this is one of those cases where it's necessary to take a step back, before taking two steps forward. After this small refactor, we can see that all of our assertions are returning an `error` value. We can use this, to create another wrapping function, which will accept a _varadiac_ `error` parameter. We will iterate over these errors, returning the first value which doesn't equal `nil`. If all errors are `nil`, we will simply return `nil`:

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

A nice side effect of the last few refactors, is that we can use _any_ logical operator, with our `validate` function, to validate our required properties. We can now validate any type we wish, using this method. Furthermore, we aren't even coupled to our `validate` function, when using our `validateMany` function. Our `validateMany` function simply takes `...error` as input, so we can use any function which returns an error, in our validation. This opens up a lot of different options, as an example, we could extract our validation of the `User` into a method:

> NOTE: It can be argued, that the function name `validateMany` is a little too specific to our usage. We could easily rename this function to something more generic, such as `handleErrors` or `returnErrorIfNotNil`. Howecer, for the purposes of this article, the current name will do just fine.

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

We now have a pretty solid solution, which is easy to implement for any of our structs with fields we wish to check as a requirement. However, this still isn't ideal. For each `struct` with a required field, we now have to implement a check on these fields. It would be much better if we could set a `required` tag on our structs. Unfortunately, despite this being proposed a few times, the Go programming language team have rejected this consequently. Their argument being, that they wish to keep the `json` package small and that it would be entirely possible to implement an equivalent feature oneself.

## The Advanced Solution
So, let's try to implement this ourselves. I'm not going to lie. This won't be pretty, but I will try my best to show an example of how to implement this, in the form of a required struct type. Let's being by showing an example of implementing a required string type. We will start out with the easy part and create a wrapper struct for our `string` value:

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
	ErrStringEmpty = errors.New("type of required.String not allowed to be empty")
	// ErrCannotUnmarshal represents an unmarshaling error
	ErrCannotUnmarshal = fmt.Errorf("json: cannot unmarshal into Go value of type required.String")
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

We aren't really doing anything crazy here. Essentially, we just use the standard `json.Unmarshal` into an `interface{}` and then checking the type of the unmarshalled JSON. If it turns out that it's a `string`, then we check if the value is equivalent of empty. If it is, then we return an error and if it's not, then we will assign the value to our `String::value` and return no error. Of course, if we find out it's not a `string` (which hopefully should never happen), we will return an error displaying our disgust.

> NOTE: We are passing an `interface{}` to the `json.Unmarshal` function, as the second parameter must be a pointer. Go does not allow us to reference a primitive type, however, we can reference a type of `interface{}`. Therefore, if we pass the pointer of the `interface{}` and then do a type conversion to `string`, we are essentially passing a pointer to our primitive `string` type. 

 Of course, we also need to be able to marshal `String` to a JSON string (and not only from). For this, we need to implement the `json.Marshaler` interface, by implementing the `MarshalJSON` method on our `String`. 

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

Running these tests, gives us the following output:
```
--- FAIL: TestStringValidation (0.00s)
    --- FAIL: TestStringValidation/nil_string (0.00s)
        /Users/lassemartinjakobsen/projects/json-validation/required/string_test.go:34: <nil>
```

Wait? What the fuck? So it turns out, that life is just not that simple in the land of Go. The reason for this, is that when we pass the `"{}"` JSON string, the `json.Unmarshal` will ignore the fields which aren't present. This means that because the the `name` property doesn't appear in the JSON, the `UnmarshalJSON` function is never called, as it is completely skipped. Therefore, we never actually get a chance to check whether our required field actually contains anything, and therefore our tests fail miserably. 

This means, that we need to find a way of checking whether our parsed required struct's values are empty. However, we don't want to go back to the validation strategy, as this would render our progress completely redundant. Therefore, we will need to find a way of doing this somewhat generically, while still using Go. You might have already seen this coming, but I still hate to say this... we are going to have to use the `reflect` package :grimacing:

Our approach will be to use our strategy from earlier in this article. First we will invoke the `json.Unmarshal` function, then our own function:

```go
// ReturnIfError will iterate over a variadac error and return
// an error if the given value is not nil
func ReturnIfError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// Unmarshal is a wrapping function of the json.Unmarshal function
func Unmarshal(data []byte, v interface{}) error {
	return ReturnIfError(
		json.Unmarshal(data, v),
		CheckValues(v),
	)
}
```

So far so good. But, now we can no longer avoid our confrontation with the `reflect` package:

```go
// CheckValues will check the values of a given interface and ensure
// that if it contains a required struct, that the required values
// are not empty
func CheckValues(v interface{}) error {
	vo := reflect.ValueOf(v)
	for vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}
	return CheckStructIsRequired(vo)
}
```

The code above is relatively simple. We are accepting an `interface{}`. We are then checking if this value is a pointer. If it is a pointer, we will call the `Elem()` function, which will return the value of the pointer. We will keep doing this, until we have an actual value. Once we have the value, we will pass it to the next function `CheckStructIsRequired`:

```go
// CheckStructIsRequired will inspect the given reflect.Value. If it contains
// a required struct, it will check it's content, if it contains a struct
// it will recursively invoke the function once more, if none of these apply
// nil will be returned.
func CheckStructIsRequired(vo reflect.Value) error {
	if vo.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < vo.NumField(); i++ {
		vtf := vo.Field(i)
		switch vtf.Type() {
		case reflect.TypeOf(String{}):
			return checkRequiredValue(vtf)
		}
		if vtf.Kind() == reflect.Struct {
			if err := CheckStructIsRequired(vtf); err != nil {
				return err
			}
		}
	}
	return nil
}
```

In this function we will iterate all the fields of our struct contained in the given `reflect.Value`. For each field, we will check if the `Type()` is equal to our required `String` type. If it is, then we will continue to check the values in the `checkRequiredValue`. If the value is not a required type, we will check whether it's a struct. If it is, we will recursively invoke our `CheckStructIsRequired`, but if it isn't we will just continue our for loop. So, the only thing we have left to look at is our `checkRequiredValue` function:

```go
func checkRequiredValue(vo reflect.Value) error {
	for i := 0; i < vo.NumField(); i++ {
		vtf := vo.Field(i)
		switch vtf.Kind() {
		case reflect.String:
			if vtf.Len() == 0 {
				return ErrStringEmpty
			}
		}
	}
	return nil
}
```
Once again, we will iterate over the properties of the given `reflect.Value` and then, using a switch statement, determine the type of our property. In our current case, we are only interested in `string` values. Also, fortunately, all we need to do, to check whether the value has been set, is by checking the length of the given value. If the value is `0`, then we can conclude that the value definitely hasn't been set. 

If we use our `UnmarshalJSON` function in our tests, instead of the `json.Unmarshal` function, we can now see our tests passing. Which is wonderful. We can also remove the check from our `UnmarshalJSON` implementation on our `String` struct, as this is now superfluous. 

Of course, this solution isn't particularly complete, but it's a good start. All we need to do to support other types, is to implement the `UnmarshalJSON` and `MarshalJSON` interface methods and then include the type in our switch statement, in the `checkRequiredValue` function. I can recommend looking into the package 
"github.com/davecgh/go-spew/spew", if you want to look at a more full implementation, using a lot of the same techniques as described here. 

Either way, the end goal of this implementation type, is to allow developers to use our package (let's call it `required`), to 'tag' struct properties as required JSON field. The structs ending up looking something like the following:

```go
type User struct {
    FirstName required.String
    LastName required.String
    JobTitle required.String
    Twitter required.String
    Stats Stats
}

type Stats struct {
    Tweets required.Int64
}
```

This seems a lot tidier, from the perspective of whomever is implementing these fields on the structs. It's a lot of code to get this point, but the ROI does become apparent, the more that you use it. 

## Summary
The last solution might seem exiting because of the usage of reflect and how it ended up being a rather generic solution. However, keep in mind, that the complexity of the code has hugely increased. The debugging ease and general type safety of the code, has also dimished significantly. This is quite a big trade-off and therefore the last solution is by no means perfect, nor suited for every situation. As said previously, the right solution for a task, depends very much on the task. If we know that we are going to use the `required` package throughout our code base with hundreds of structs. Then it's probably appropriate to development and implement, however, if we are only talking about ensuring requirements for a few fields. Then the more manual method of checking the required fields is much more appropriate. 

I hope that this article gave you some insight into some different techniques of approaching an issue, using Go and all feedback is extremely welcome. Thank you so much for spending your time reading this. If you want to keep up-to-date on my articles, be sure to follow me on Twitter @ifndef_lmj. Thanks again!

