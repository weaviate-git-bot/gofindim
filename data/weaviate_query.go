package data

import (
	"context"
	"fmt"
	"math"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"gocv.io/x/gocv"
)

func descriptorToFloatSlice(desc *gocv.Mat) []float32 {
	floatSlice := make([]float32, desc.Total())
	for i := 0; i < len(floatSlice); i++ {
		value := float64(desc.GetVecfAt(0, i)[0])
		if math.IsNaN(value) || math.IsInf(value, 0) {
			// Handle NaN or infinite values, e.g., by setting a default value or ignoring them
			value = 0 // example: set to 0}
			floatSlice[i] = float32(value)
			// fmt.Printf("%v\n", floatSlice)
		}
		floatSlice[i] = float32(value)
	}
	return floatSlice
}

// func MatToFloat32(mat *gocv.Mat) []float32 {
//     data := make([]float32, mat.Rows()*mat.Cols())
//     for i, row := range mat.GetVecfAt().GetFloatAt(0, 0) {
//         data[i] = float32(row)
//     }
//     return data
// }

func matToFloat32Array(mat gocv.Mat) []float32 {
	data := make([]float32, mat.Total())
	rows := mat.Rows()
	cols := mat.Cols()

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			value := mat.GetFloatAt(row, col)
			if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
				value = 0
			}
			data[row*cols+col] = value
		}
	}
	return data
	// norm := normalizeVector(data)

	// return norm
}
func normalizeVector(vector []float32) []float32 {
	// Calculate the L2 norm (also known as the Euclidean norm or magnitude)
	var l2Norm float32 = 0
	for _, value := range vector {
		l2Norm += value * value
	}
	l2Norm = float32(math.Sqrt(float64(l2Norm)))

	// Normalize the vector if the L2 norm isn't 0
	if l2Norm != 0 {
		for i := range vector {
			vector[i] /= l2Norm
		}
	}

	return vector
}

func SearchWeaviate(i *ImageFile, certainty float32, limit int, client *weaviate.Client) error {
	fieldToQuery := "path"
	queryField := graphql.Field{Name: fieldToQuery}
	vector, err := i.ToVector()
	if err != nil {
		return err
	}
	nearVector := client.GraphQL().NearVectorArgBuilder().WithVector(vector).WithCertainty(certainty)
	response, err := client.GraphQL().Get().
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

type Image struct {
	Name string
}

type Get struct {
	Image []Image
}

type GraphQLResponse struct {
	Data map[string]Get
}

func SearchWeaviateWithVector(vector []float32, certainty float32, limit int, client *weaviate.Client) error {
	nameField := graphql.Field{Name: "path"}
	nearVector := client.GraphQL().NearVectorArgBuilder().WithVector(vector).WithCertainty(certainty)
	response, err := client.GraphQL().Get().
		WithClassName("Image").
		WithFields(nameField).
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
				fmt.Printf("%v\n", image.(map[string]interface{})["path"])
			}
		}
	}
	return nil

}
func SearchWeaviateObj(i *ImageFile) {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	orb := gocv.NewORBWithParams(1000, 1.2, 8, 31, 0, 2, gocv.ORBScoreTypeHarris, 31, 20)
	defer orb.Close()

	// _, descriptor := i.toKeypointsDescriptors(&orb)

	// sourceDescriptor := matToFloat32Array(descriptor)
	// nearVector := client.GraphQL().NearVectorArgBuilder().WithVector(sourceDescriptor).WithCertainty(0.1)
	nearObject := client.GraphQL().NearObjectArgBuilder().WithID("413910a2-1bd3-45e6-8c57-4a02e1798ffa")
	// floatsAsJSON := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(sourceDescriptor)), ","), "[]")
	imageDataField := graphql.Field{Name: "imageData"}
	// _additional := graphql.Field{
	// 	Name: "_additional", Fields: []graphql.Field{
	// 		{Name: "certainty"}, // only supported if distance==cosine
	// 		{Name: "distance"},  // always supported
	// 	},
	// }
	response, err := client.GraphQL().Get().
		WithClassName("Image").
		WithFields(imageDataField).
		WithNearObject(nearObject).
		WithLimit(10).
		Do(context.Background())
	// Check error and handle response
	if err != nil {
		panic(err)
	}
	match := response
	fmt.Printf("%v\n", match.Data)

}
