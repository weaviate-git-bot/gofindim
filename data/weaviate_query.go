package data

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

func SearchWeaviate(i *ImageFile, certainty float32, limit int, client *weaviate.Client) error {
	fieldToQuery := "name"
	vector, err := i.ToVector()
	queryField := graphql.Field{Name: fieldToQuery}
	if err != nil {
		return err
	}
	nearVector := client.GraphQL().NearVectorArgBuilder().WithVector(vector).WithCertainty(certainty)
	qb := NewQueryBuilder(client)
	response, err := qb.
		WithClassName("Image").
		WithFields(queryField).
		WithNearVector(nearVector).
		WithLimit(limit).
		Do(context.Background())
	// Check error and handle response
	if err != nil {
		return err
	}
	for _, outerMap := range response.Data {
		for _, images := range outerMap.(map[string]interface{}) {
			for _, image := range images.([]interface{}) {
				fmt.Printf("%v\n", image.(map[string]interface{})[fieldToQuery])
			}
		}
	}
	return nil
}

func SearchWeaviateWithVector(vector []float32, certainty float32, limit int, fields []string, client *weaviate.Client) error {
	nearVector := client.GraphQL().NearVectorArgBuilder().WithVector(vector).WithCertainty(certainty)
	response, err := client.GraphQL().Get().
		WithClassName("Image").
		WithNearVector(nearVector).
		WithLimit(limit).
		Do(context.Background())
	// Check error and handle response
	if err != nil {
		return err
	}
	// var outerMap map[string]interface{}
	data := response.Data
	results, err := ParseImageData(data, fields)
	if err != nil {
		return err
	}
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("error marshalling results: %v", err)
	}
	fmt.Printf("%v\n", string(resultsJSON))

	// outerValue := data["Get"]
	// println(data["Get"].(map[string]interface{})["Image"].([]interface{})[0].(map[string]interface{})["path"])
	// if outerMap, ok := outerValue.(map[string]interface{}); ok {
	// 	if images, ok := outerMap["Image"].([]interface{}); ok {
	// 		if image, ok := images[0].(map[string]interface{}); ok {
	// 			if path, ok := image["path"]; ok {
	// 				fmt.Println(path)
	// 			}
	// 		}
	// 	}
	// } else {
	// 	return fmt.Errorf("outerValue is not a map[string]interface{}")
	// 	// handle error
	// }
	// err := json.Unmarshal([]byte(data.(string)), &outerMap)
	// err := json.Unmarshal((outerMap.(map[string]interface{})["Get"]), &outerMapData)
	// fmt.Printf("%v\n", outerMapData)
	// if err != nil {
	// 	return err
	// }
	return nil
	// for _, outerMap := range response.Data {
	// 	for _, images := range outerMap.(map[string]interface{}) {
	// 		for _, image := range images.([]interface{}) {
	// 			fmt.Printf("%v\n", image.(map[string]interface{})["path"])
	// 		}
	// 	}
	// }
	// return nil
}

func SearchWeaviateWithImageFile(i *ImageFile, certainty float32, limit int, fields []string, client *weaviate.Client) (*[]ImageNode, error) {
	qb := NewQueryBuilder(client)
	response, err := qb.
		NearImage(i, certainty).
		SelectFields(fields).
		WithClassName("Image").
		WithLimit(limit).
		Do(context.Background())
	// Check error and handle response
	if err != nil {
		return nil, err
	}
	results, err := ParseImageData(response.Data, fields)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
func SearchWeaviateWithText(t string, certainty float32, limit int, fields []string, client *weaviate.Client) (*[]ImageNode, error) {
	qb := NewQueryBuilder(client)
	response, err := qb.
		NearText(t, certainty).
		SelectFields(fields).
		WithClassName("Image").
		WithLimit(limit).
		Do(context.Background())
	// Check error and handle response
	if err != nil {
		return nil, err
	}
	results, err := ParseImageData(response.Data, fields)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
