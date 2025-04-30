package utils

import (
	"encoding/json"

	"gorm.io/datatypes"
)


func EncodeJSONSlice(data interface{}) (datatypes.JSON, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(jsonData), nil
}

// DecodeJSON converts GORM datatypes.JSON to a Go slice
func DecodeJSONSlice(jsonData datatypes.JSON, result interface{}) error {
	if len(jsonData) == 0 || string(jsonData) == "null" {
		return nil
	}

	err := json.Unmarshal(jsonData, result)
	if err != nil {
		return err
	}
	return nil
}
