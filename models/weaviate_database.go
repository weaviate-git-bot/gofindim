package models

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"gocv.io/x/gocv"
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

func keypointsToBytes(keypoints []gocv.KeyPoint) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(keypoints)
	return buf.Bytes(), err
}

func matToBytes(mat gocv.Mat) ([]byte, error) {
	data := mat.ToBytes()
	return data, nil
}

func generateUID() uuid.UUID {
	// Generate a UUID
	newUUID := uuid.New()
	// Print the UUID
	return newUUID

}

func InsertIntoWeaviate(img *ImageFile) error {
	// orb := gocv.NewORBWithParams(1000, 1.2, 8, 31, 0, 2, gocv.ORBScoreTypeHarris, 31, 20)
	// defer orb.Close()
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// dataUid := generateUID()
	imgUid := generateUID()
	fmt.Println("imgUid: ", imgUid)

	// keypoints, descriptors := img.toKeypointsDescriptors(&orb)
	imageInterface := map[string]interface{}{
		"name":  img.Name,
		"path":  img.Path,
		"image": img.Base64,
	}
	// imageDataInterface := map[string]interface{}{
	// 	"filename": img.Name,
	// 	"rating":   5, // sample rating
	// }

	// err = client.Data().Validator().
	// 	WithID(dataUid.String()).
	// 	WithClassName("ImageData").
	// 	WithProperties(imageDataInterface).
	// 	Do(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
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
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// _, err = client.Batch().ObjectsBatcher().
	// 	WithObjects(imageDataObject, imageObject).
	// 	WithConsistencyLevel(replication.ConsistencyLevel.ALL).
	// 	Do(context.Background())

	// if err != nil {
	// 	fmt.Errorf("Batching failed: %v", err)
	// }
	fmt.Printf("Created \n")
	return err
}

// func insertImageData(filename string, rating int) (string, error) {
// 	// Create a JSON object for ImageData
// 	data := map[string]interface{}{
// 		"class": "ImageData",
// 		"properties": map[string]interface{}{
// 			"filename": filename,
// 			"rating":   rating,
// 		},
// 	}
// 	// POST the object and get the response
// 	// ...

// 	return idOfInsertedObject, err
// }

// func insertImage(keypoints, descriptor []byte, imageDataRef string) (string, error) {
// 	// Create a JSON object for Image
// 	data := map[string]interface{}{
// 		"class": "Image",
// 		"properties": map[string]interface{}{
// 			"keypoints":  keypoints,
// 			"descriptor": descriptor,
// 			"imageData": map[string]string{
// 				"beacon": "weaviate://localhost/ImageData/" + imageDataRef,
// 			},
// 		},
// 	}
// 	reponse, err := postToWeaviate("/objects", data)
// 	if err != nil {
// 		return "", err
// 	}
// 	fmt.Println(reponse)
// 	return reponse["id"].(string), nil

// }

// func postToWeaviate(endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("POST", "http://localhost:8080/v1"+endpoint, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	responseBody, err := ioutil.ReadAll(resp.Body)

// 	var responseMap map[string]interface{}
// 	json.Unmarshal(responseBody, &responseMap)

// 	return responseMap, err
// }

// func keypointsToJSON(kp []gocv.KeyPoint, desc gocv.Mat, imageID string) []map[string]interface{} {
// 	var jsonData []map[string]interface{}

// 	for i, keypoint := range kp {
// 		x, y := keypoint.X, keypoint.Y
// 		coordinates := []float64{x, y}
// 		descriptor := desc.ToBytes()
// 		descriptorStr := base64.StdEncoding.EncodeToString(descriptor)

// 		keypointJSON := map[string]interface{}{
// 			"class": "Keypoint",
// 			"properties": map[string]interface{}{
// 				"coordinates": coordinates,
// 				"descriptor":  descriptorStr,
// 				"image":       map[string]string{"beacon": "weaviate://localhost/" + imageID},
// 			},
// 		}
// 		jsonData = append(jsonData, keypointJSON)
// 	}
// 	return jsonData
// }
