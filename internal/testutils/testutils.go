package testutils

import (
	"errors"
	"reflect"
)

// ExampleData serves as sample data for the test methods
type ExampleData struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

// Unmarshal a map of interfaces values to ExampleData's fields
func (ed *ExampleData) Unmarshal(data map[string]interface{}) error {
	if val, ok := data["message"].(string); ok == true {
		ed.Message = val
	} else {
		return errors.New("Failed to get the message")
	}
	if val := int(data["count"].(float64)); reflect.TypeOf(val).Kind() == reflect.Int {
		ed.Count = val
	} else {
		return errors.New("Failed to get the count")
	}

	return nil
}

// UnmarshalFromInterface unmarshals an interface which is a map value to ExampleData's fields
func (ed *ExampleData) UnmarshalFromInterface(data interface{}) error {

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Map {
		m := data.(map[string]interface{})
		if val, ok := m["message"].(string); ok == true {
			ed.Message = val
		} else {
			return errors.New("Failed to unmarshal field message")
		}
		if val := int(m["count"].(float64)); reflect.TypeOf(val).Kind() == reflect.Int {
			ed.Count = val
		} else {
			return errors.New("Failed to unmarshal field count")
		}
	} else {
		return errors.New("The data type is not a map of type map[string]interface{}. Got: " + v.Kind().String())
	}

	return nil
}
