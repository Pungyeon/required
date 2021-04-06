# Required JSON tags

[![GoDoc](https://godoc.org/github.com/Pungyeon/required?status.svg)](https://godoc.org/github.com/Pungyeon/required) 
[![Go Report Card](https://goreportcard.com/badge/github.com/Pungyeon/required)](https://goreportcard.com/report/github.com/Pungyeon/required)  

## Introduction
Long story short, a junior engineer at work asked me (a long time ago now). How to make `JSON` fields required. I didn't have a good answer for him, and after some short research it turns out that there isn't a good answer for this!

I therefore wrote the following article on making a package validating `JSON` input: [article/README.md](Original Article).

However, a delve into the rabbit hole later, I have made my own `JSON` parser, which does support required tags. It is definitely not a finished or production ready product, so all help and feedback is extremely welcome. The original article was meant to solve a problem, whereas the `JSON` parser spawned from pure curiosity.

## Usage
It should be possible to substitute the std `json` import with the current library, without issue.

If this is not the case, please make sure to report this as an issue here on GitHub :)

### JSON Unmarshalling
The usage is very simple. Adding the `JSON` struct tag "required", will have the parser enforce this field to be present, when unmarshalling:

```go
type User struct {
  FirstName string `json:"first_name,required"`
  LastName  string `json:"last_name,required"`
  Email     string `json:"email,required"`
  GitHub    string `json:"github"`
  LinkedIn  string
}
``` 

In the above example, `FirstName`, `LastName` and `Email` are required, where as `GitHub` and `LinkedIn` are not.

```go
func main() {
    var user User
    if err := Unmarshal([]byte(`{}`), &user); err != nil {
        panic(err) // this will return an error because of the missing required fields
    }
    fmt.Println(user)
}
```

Where as the following JSON string will parse without returning an error:

```go
`{
  "first_name": "lasse",
  "last_name": "jakobsen",
  "email": "lasse@jakobsen.dev"
}`
```

Furthermore, it is possible to implement the `Required` interface, to create custom validation of a type.

```go
type Required interface {
	IsValueValid() error
}
```

Below is an example of implementing this interface for a type alias of `string` called `Email`. This will validate the string to adhere to the general structure of an e-mail. This means, that when being parsed, the `json` will be ensured to not only exist, but also that it's valid as defined in the `IsValidValue` function.

```go
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
```

Refer to samples for a more detailed example of this.

### Marshalling
As of writing this document, this library is currently using a custom `json.Marshal` and `json.Encoder`. This library *does not currently support `required` tag checking*, please show your interest, if you would like this by creating a new issue. The `json.Marshal` function is compatible with the standard library functionality. Though, substantially faster:

```
goos: darwin
goarch: amd64
pkg: github.com/Pungyeon/required/pkg/json
BenchmarkMarshalStd-8             603121              1859 ns/op             752 B/op         12 allocs/op
BenchmarkMarshalPkg-8             729928              1581 ns/op             832 B/op         10 allocs/op
PASS
ok      github.com/Pungyeon/required/pkg/json   2.445s
```




