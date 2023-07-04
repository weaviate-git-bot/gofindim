package data

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/agentx3/gofindim/utils"
	"github.com/weaviate/weaviate/entities/models"
)

func ParseImageData(data map[string]models.JSONObject, fields []string) ([]ImageNode, error) {
	results := make([]ImageNode, 0)
	if len(data) == 0 {
		return results, nil
	}

	outerValue, ok := data["Get"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error parsing outerValue: expected map[string]interface{}, but got %T", data["Get"])
	}

	images, ok := outerValue["Image"].([]interface{})
	if !ok {
		for key, value := range outerValue {
			fmt.Println(key, value)
		}
		return nil, fmt.Errorf("error parsing outerMap: expected []interface{}, but got %T", outerValue["Image"])
	}

	for _, image := range images {
		img, ok := image.(map[string]interface{})
		if !ok {
			return results, nil
		}
		node := &ImageNode{}
		nodeValue := reflect.ValueOf(node).Elem()
		nodeType := reflect.TypeOf(ImageNode{})

		for i := 0; i < nodeValue.NumField(); i++ {
			field := nodeType.Field(i).Name
			key := strings.ToLower(field)

			if utils.StringInSlice(key, []string{"id", "distance"}) {
				val, err := additionalFieldValid(key, img)
				if err == nil {
					nodeValue.Field(i).Set(reflect.ValueOf(val))
				}

			} else if value, ok := img[key]; ok {
				nodeValue.Field(i).Set(reflect.ValueOf(value))
			}
		}
		results = append(results, *node)
	}

	return results, nil
}

func additionalFieldValid(field string, img map[string]interface{}) (interface{}, error) {
	additionals, ok := img["_additional"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error parsing %v: expected map[string]interface{}, but got %T", field, img["Image"])
	}
	fieldVal := additionals[field]
	if fieldVal == nil {
		return nil, fmt.Errorf("error parsing additionals[\"%v\"]: expected string, but got %T", field, additionals[field])
	}
	return fieldVal, nil
}
