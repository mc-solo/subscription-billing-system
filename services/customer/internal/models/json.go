package models

import (
	"encoding/json"
	"fmt"
)

type JSON map[string]interface{}

func (j JSON) Value() (interface{}, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSON)
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("cannot scan type of %T into JSON", value)
	}

	return json.Unmarshal(data, j)
}
