package data

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/weaviate/weaviate/entities/models"
)

func ParseImageData(data map[string]models.JSONObject, fields []string) ([]ImageNode, error) {
	results := make([]ImageNode, 0)
	if len(data) == 0 {
		return results, nil
	}
	outerValue := data["Get"]
	// println(data["Get"].(map[string]interface{})["Image"].([]interface{})[0].(map[string]interface{})["path"])
	if outerMap, ok := outerValue.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("error parsing outerValue: %v", outerValue)
	} else if images, ok := outerMap["Image"].([]interface{}); !ok {
		return nil, fmt.Errorf("error parsing outerMap: %v", outerMap)
	} else {
		for _, image := range images {
			if img, ok := image.(map[string]interface{}); !ok {
				return results, nil
			} else {
				node := &ImageNode{}
				nodeValue := reflect.ValueOf(node).Elem()
				nodeType := reflect.TypeOf(ImageNode{})
				for i := 0; i < nodeValue.NumField(); i++ {
					field := nodeType.Field(i).Name
					key := strings.ToLower(field)
					if value, ok := img[key]; ok {
						nodeValue.Field(i).Set(reflect.ValueOf(value))
					}
				}
				results = append(results, *node)
			}
		}
	}

	return results, nil
}
