package util

import (
	"encoding/json"
	"fmt"
)

// Pformat returns a pretty format output of any value.
func Pformat(value interface{}) (string) {
	if s, ok := value.(string); ok {
		return s
	}
	valueJson, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("Could not marshal the value; error %v", err)
	}
	return string(valueJson)
}
