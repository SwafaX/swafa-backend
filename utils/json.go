package utils

import (
	"encoding/json"
)

// ParseJSON parses a JSON string into a struct
func ParseJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// ToJSON converts a struct to a JSON string
func ToJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
} 