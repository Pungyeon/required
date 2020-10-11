### Required JSON tag
Long story short, a junior engineer at work asked me (a long time ago now). How to make `JSON` fields required. I didn't have a good answer for him, and after some short research it turns out that there isn't a good answer for this!

I therefore wrote the following article on making a package validating `JSON` input: [article/README.md](Original Article).

However, a delve into the rabbit hole later, I have made my own `JSON` parser, which does support required tags. It is definitely not a finished or production ready product, so all help and feedback is extremely welcome. The original article was meant to solve a problem, whereas the `JSON` parser spawned from pure curiosity.

### Usage
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

Refer to samples for a more detailed example of this.