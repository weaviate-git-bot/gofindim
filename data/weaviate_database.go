package data

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

type ImageDataNode struct {
	Filename string `json:"filename"`
	Rating   int    `json:"rating"`
}

type ImageNode struct {
	Coordinates string        `json:"coordinates"`
	Descriptor  []float64     `json:"descriptor"`
	Image       ImageDataNode `json:"image"`
}

func generateUID() uuid.UUID {
	// Generate a UUID
	newUUID := uuid.New()
	// Print the UUID
	return newUUID

}

func CreateImageClass(client *weaviate.Client) error {
	// Create a class
	err := client.
		Schema().
		ClassCreator().
		WithClass(ImageClass).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func InsertIntoWeaviate(img *ImageFile, client *weaviate.Client) error {
	var err error
	// dataUid := generateUID()
	imgUid := generateUID()
	fmt.Println("imgUid: ", imgUid)

	vector, err := img.ToVector()
	if err != nil {
		return err
	}
	imageInterface := map[string]interface{}{
		"name":      img.Name,
		"path":      img.Path,
		"embedding": vector,
	}
	// imageDataInterface := map[string]interface{}{
	// 	"filename": img.Name,
	// 	"rating":   5, // sample rating
	// }

	err = client.Data().Validator().
		WithID(imgUid.String()).
		WithClassName("Image").
		WithProperties(imageInterface).
		Do(context.Background())
	if err != nil {
		fmt.Printf("There's been an Error: %v\n", err)
		return err
	}
	_, err = client.Data().Creator().
		WithClassName("Image").
		WithID(imgUid.String()).
		WithProperties(imageInterface).
		WithVector(vector).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created \n")
	return err
}

func InsertMultipleIntoWeaviate(img []*ImageFile, client *weaviate.Client) error {
	var err error
	// imgUid := generateUID()
	// fmt.Println("imgUid: ", imgUid)

	var objects = make([]*models.Object, 0, len(img))
	for _, image := range img {
		vector, err := image.ToVector()
		if err != nil {
			return fmt.Errorf("error converting image to vector: %v", err)
		}
		object := models.Object{
			Class: "Image",
			Properties: map[string]interface{}{
				"name":      image.Name,
				"path":      image.Path,
				"embedding": vector,
			},
			Vector: vector,
		}
		objects = append(objects, &object)
	}

	batcher := client.Batch().ObjectsBatcher()
	for _, object := range objects {
		batcher.WithObject(object)
	}
	_, err = batcher.Do(context.Background())
	if err != nil {
		return fmt.Errorf("error inserting batch: %v", err)
	}
	fmt.Printf("Created \n")
	return err
}
