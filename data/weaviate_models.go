package data

import "github.com/weaviate/weaviate/entities/models"

var ImageClass = &models.Class{
	Class:           "Image",
	Vectorizer:      "none", // If set to "none" you must always provide vectors yourself. Could be any other "text2vec-*" also.
	VectorIndexType: "hnsw",
	ModuleConfig: map[string]interface{}{
		"multi2vec-clip": map[string]interface{}{
			"model": "sentence-transformers-clip-ViT-B-32",
			"options": map[string]interface{}{
				"waitForModel": false,
			},
			// "imageFields": []string{"image"},
		},
	},
	Properties: []*models.Property{
		&models.Property{
			Name:     "name",
			DataType: []string{"text"},
		},
		&models.Property{
			Name:        "path",
			Description: "Path to the image file",
			DataType:    []string{"text"},
		},
		&models.Property{
			Name:        "embedding",
			Description: "The image itself. This is a base64 encoded image.",
			DataType:    []string{"number[]"},
		},
		// &models.Property{
		// 	Name:     "imageData",
		// 	Description: "The metadata of the image.",
		// 	DataType: []string{"ImageData"},
		// },
	},
}

// var ImageDataClass = &models.Class{
// 	Class:           "ImageData",
// 	Vectorizer:      "multi2vec-clip", // If set to "none" you must always provide vectors yourself. Could be any other "text2vec-*" also.
// 	VectorIndexType: "hnsw",
// 	ModuleConfig: map[string]interface{}{
// 		"multi2vec-clip": map[string]interface{}{
// 			"model": "sentence-transformers-clip-ViT-B-32",
// 			"options": map[string]interface{}{
// 				"waitForModel": false,
// 			},
// 			"textFields": []string{"tags"},
// 		},
// 	},
// 	Properties: []*models.Property{
// 		&models.Property{
// 			Name:     "author",
// 			DataType: []string{"text"},
// 		},
// 		&models.Property{
// 			Name:     "path",
// 			Description:     "Path to the image file",
// 			DataType: []string{"text"},
// 		},
// 		&models.Property{
// 			Name:     "tags",
// 			Description: "Tags of the image.",
// 			DataType: []string{"string"},
// 		},
// 	},
// }
