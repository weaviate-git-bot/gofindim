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
		},
	},
	Properties: []*models.Property{
		{Name: "name",
			DataType: []string{"text"},
		},
		{Name: "path",
			Description: "Path to the image file",
			DataType:    []string{"text"},
		},
	},
}
