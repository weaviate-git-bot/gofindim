package data

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/agentx3/gofindim/utils"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

type ImageNode struct {
	Name   string `json:"name,omitempty"`
	Id     string `json:"id,omitempty"`
	Rating int    `json:"rating,omitempty"`
	Path   string `json:"path,omitempty"`
	// This should almost always be a float32, but weaviate returns it as a float64
	// so we have to use a float64 here. It's important to note that this should
	// only be used to hold values for reading from the database.
	Distance float64 `json:"distance,omitempty"`
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
	imgUid, err := FileToUUID(img.Path)
	if err != nil {
		return err
	}
	fmt.Println("imgUid: ", imgUid)

	vector, err := img.ToVector()
	if err != nil {
		return err
	}

	err = client.Data().Validator().
		WithID(imgUid).
		WithClassName("Image").
		WithProperties(img.toInterface()).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't valid image during insertion.", err)
	}
	_, err = client.Data().Creator().
		WithClassName("Image").
		WithID(imgUid).
		WithProperties(img.toInterface()).
		WithVector(vector).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't insert image into database.", err)
	}
	fmt.Printf("Created \n")
	return err
}

func InsertMultipleIntoWeaviate(img []ImageFile, client *weaviate.Client) error {
	var err error
	var objects = make([]*models.Object, 0, len(img))
	for _, image := range img {
		vector, err := image.ToVector()
		if err != nil {
			return fmt.Errorf("error converting image to vector: %v", err)
		}
		uuid, err := FileToUUID(image.Path)
		if err != nil {
			return err
		}
		object := models.Object{
			Class: "Image",
			ID:    strfmt.UUID(uuid),
			Properties: map[string]interface{}{
				"name": image.Name,
				"path": image.Path,
				// "embedding": vector,
			},
			Vector: vector,
		}
		err = client.Data().Validator().
			WithID(uuid).
			WithClassName("Image").
			WithProperties(image.toInterface()).
			Do(context.Background())
		if err != nil {
			fmt.Printf("Couldn't validate image during insertion of %v : %v", image.Path, err)
			continue
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
	fmt.Printf("Created %v\n objects.", len(objects))
	return err
}

func InsertDirectoryIntoWeaviate(path string, client *weaviate.Client) error {
	var err error
	imagePaths := []string{}
	// images := []*models.ImageFile{}

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if strings.HasPrefix(d.Name(), ".") && d.IsDir() {
			fmt.Println("Skipping directory", path)
			return filepath.SkipDir
		}
		if utils.IsImage(path) {
			imagePaths = append(imagePaths, path)
			// img, err := models.NewImageFileFromPath(path)
			// if err != nil {
			// 	return fmt.Errorf("error creating image file from path: %w", err)
			// }
			// println("Adding", len(images), "images")
			// images = append(images, img)
		}
		return nil
	})
	batchSize := 50
	waitGroup := sync.WaitGroup{}

	for i := 0; i < len(imagePaths); i += batchSize {
		end := i + batchSize
		imageBatch := []ImageFile{}
		if end > len(imagePaths) {
			end = len(imagePaths)
		}
		batch := imagePaths[i:end]
		for _, img := range batch {
			imageFile, err := NewImageFileFromPath(img)
			if err != nil {
				fmt.Printf("error creating image file from path %v: %v", img, err)
				continue
			}
			imageBatch = append(imageBatch, *imageFile)
		}
		waitGroup.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err = InsertMultipleIntoWeaviate(imageBatch, client)
			if err != nil {
				fmt.Printf("error inserting multiple into weaviate: %v", err)
			}
		}(&waitGroup)
		fmt.Printf("Processing batch: %v\n", i/batchSize)
	}
	waitGroup.Wait()
	// err := models.InsertMultipleIntoWeaviate(images, client)
	// if err != nil {
	// 	return fmt.Errorf("error inserting multiple into weaviate: %w", err)
	// }
	fmt.Println("Inserted multiple into weaviate")
	return nil
}

func VectorFromUUID(id string, client *weaviate.Client) ([]float32, error) {

	// Open the file
	response, err := client.Data().
		ObjectsGetter().
		WithVector().
		WithClassName("Image").
		WithID(id).
		Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting object with id %s: %s\n", id, err)
		return nil, err
	}
	if len(response) == 0 {
		return nil, fmt.Errorf("No object with id %s found", id)
	}
	result := *response[0]
	return result.Vector, nil

	// embeddingInterfaces := result.Properties.(map[string]interface{})["embedding"].([]interface{})
	// embeddingValues := make([]float32, len(embeddingInterfaces))

	// for i, val := range embeddingInterfaces {
	// 	floatVal, ok := val.(float64)
	// 	if !ok {
	// 		fmt.Println("Error converting to float32")
	// 		break
	// 	}
	// 	embeddingValues[i] = float32(floatVal)
	// }

}
