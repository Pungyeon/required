package json

import stdjson "encoding/json"

func Marshal(v interface{}) ([]byte, error) {
	return stdjson.Marshal(v)
}
