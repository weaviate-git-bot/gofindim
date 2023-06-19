package data

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/agentx3/gofindim/utils"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

func SearchWeaviate(i *ImageFile, distance float32, limit int, client *weaviate.Client) error {
	fieldToQuery := "name"
	vector, err := i.ToVector()
	queryField := graphql.Field{Name: fieldToQuery}
	if err != nil {
		return err
	}
	nearVector := client.GraphQL().NearVectorArgBuilder().WithVector(vector).WithDistance(distance)
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

func SearchWeaviateWithVector(vector []float32, distance float32, limit int, fields []string, client *weaviate.Client) (*[]ImageNode, error) {
	qb := NewQueryBuilder(client)
	response, err := qb.
		NearVector(vector, distance).
		SelectFields(fields).
		WithClassName("Image").
		WithLimit(limit).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	// var outerMap map[string]interface{}
	results, err := ParseImageData(response.Data, fields)
	println(len(results))
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func SearchWeaviateWithImagePath(path string, distance float32, limit int, fields []string, client *weaviate.Client) (*[]ImageNode, error) {
	imgFile, err := NewImageFileFromPath(path)
	if err != nil {
		return nil, err
	}
	results, err := SearchWeaviateWithImageFile(imgFile, distance, limit, fields, client)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func SearchWeaviateWithImageFile(i *ImageFile, distance float32, limit int, fields []string, client *weaviate.Client) (*[]ImageNode, error) {
	qb := NewQueryBuilder(client)
	response, err := qb.
		NearImage(i, distance).
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
func SearchWeaviateWithText(t string, distance float32, limit int, fields []string, client *weaviate.Client) (*[]ImageNode, error) {
	qb := NewQueryBuilder(client)
	fmt.Printf("Searching for %s\n", fields)
	response, err := qb.
		NearText(t, distance).
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

func SearchWeaviateWithFormFile(file *multipart.FileHeader, distance float32, limit int, fields []string, client *weaviate.Client) (*[]ImageNode, error) {
	// Open the file
	f, err := file.Open()
	defer f.Close()
	if err != nil {
		return nil, err
	}
	imageFile := NewImageFileFromFormFile(f, file.Filename)
	results, err := SearchWeaviateWithImageFile(imageFile, distance, limit, []string{"path", "name"}, client)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func SearchWeaviateWithUUID(id string, distance float32, limit int, fields []string, client *weaviate.Client) (*[]ImageNode, error) {
	// Open the file
	response, err := client.Data().
		ObjectsGetter().
		WithClassName("Image").
		WithID(id).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(response) == 0 {
		return nil, fmt.Errorf("No object with id %s found", id)
	}
	result := *response[0]
	results, err := SearchWeaviateWithVector(result.Vector, distance, limit, fields, client)
	// results, err := ParseImageData(response, fields)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func SearchWeaviateWithTextAndImage(text string, image *ImageFile, textWeight, imageWeight, distance float32, limit int, fields []string, client *weaviate.Client) (*[]ImageNode, error) {
	qb := NewQueryBuilder(client)
	fmt.Printf("Searching for %s\n", fields)
	textVec, imgVec, err := VectorizeTextAndImage(text, image)
	if err != nil {
		return nil, err
	}
	combinedVec := utils.AverageVectors(textVec, imgVec, textWeight, imageWeight)

	response, err := qb.
		NearVector(combinedVec, distance).
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
