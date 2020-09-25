package json

type ObjectType int

const (
	Unknown = 0
	String  = 1
	Integer = 2
	Float   = 3
	Slice   = 4
	Obj     = 5
)

type Object struct {
	Value interface{}
	Type  ObjectType
}

func (obj *Object) add(token Token) {
	if token.Type == StringToken {
		obj.Type = String
	}
	if token.Type == FullStopToken {
		obj.Type = Float
	}
	if obj.Value == nil {
		obj.Value = token.Value
	} else {
		obj.Value = obj.Value.(string) + token.Value
	}
}
