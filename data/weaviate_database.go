package data

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

type ImageNode struct {
	Name   string `json:"name,omitempty"`
	Rating int    `json:"rating,omitempty"`
	Path   string `json:"path,omitempty"`
}

// type ImageNode struct {
// 	Coordinates string        `json:"coordinates"`
// 	Descriptor  []float64     `json:"descriptor"`
// 	Image       ImageDataNode `json:"image"`
// }

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

	err = client.Data().Validator().
		WithID(imgUid.String()).
		WithClassName("Image").
		WithProperties(img.toInterface()).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't valid image during insertion.", err)
	}
	_, err = client.Data().Creator().
		WithClassName("Image").
		WithID(imgUid.String()).
		WithProperties(img.toInterface()).
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
	_, err = client.Batch().
		ObjectsBatcher().
		WithObjects(objects...).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("error inserting batch: %v", err)
	}
	fmt.Printf("Created \n")
	return err
}
