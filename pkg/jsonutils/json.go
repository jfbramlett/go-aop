package jsonutils

import (
	"bytes"
	"encoding/json"
	"strings"
)

// ToJSON Converts the content of our msg in to a string
func ToJSON(msg interface{}) (string, error) {
	content := &strings.Builder{}
	enc := json.NewEncoder(content)
	enc.SetIndent("", "    ")
	if err := enc.Encode(msg); err != nil {
		return "", err
	}

	return content.String(), nil
}

// FromJSON converts a json string in to an object
func FromJSON(jsonContent string, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader([]byte(jsonContent)))
	err := decoder.Decode(v)
	if err != nil {
		return err
	}
	return nil
}