package utils

import (
	"encoding/json"
	"log"
	"reflect"
)

func NewPtr[T any](val T) *T {
	return &val
}

func GetValFromPtr[T any](val *T) T {
	var ret T
	if val == nil {
		return ret
	}
	return *val
}

func StructToMap(obj interface{}) map[string]interface{} {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		log.Fatalf("Error marshalling struct to JSON: %v", err)
		return nil
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON to map: %v", err)
		return nil
	}

	return result
}

func ContainsSameValues(slice []interface{}) bool {
	if len(slice) == 0 {
		return true
	}
	firstValue := slice[0]
	for _, value := range slice {
		if !reflect.DeepEqual(value, firstValue) {
			return false
		}
	}
	return true
}
