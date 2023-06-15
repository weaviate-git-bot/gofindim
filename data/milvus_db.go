package data

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func GetMilvusClient() (client.Client, error) {
	mClient, err := client.NewGrpcClient(context.Background(), "localhost:19530")
	if err != nil {
		return nil, err
	}
	return mClient, nil
}

func CreateDatabase(mClient client.Client) error {
	// create database
	err := mClient.CreateDatabase(context.Background(), databaseName)
	if err != nil {
		return err
	}
	return nil
}

func CreateCollection(mClient client.Client) error {
	// create collection
	err := mClient.CreateCollection(
		context.Background(),
		ImageSchema,
		2,
	)
	if err != nil {
		fmt.Errorf("Error creating collection", err)
	}
	if err != nil {
		return err
	}
	return nil
}

func CheckCollection(mClient client.Client, name string) (bool, error) {
	// check collection
	has, err := mClient.HasCollection(context.Background(), name)
	if err != nil {
		return false, err
	}
	return has, nil
}

func ViewCollections(mClient client.Client) error {
	// view collections
	collections, err := mClient.ListCollections(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Collections:", collections)
	return nil
}

func InitMilvus(mClient client.Client) error {
	var err error

	if err != nil {
		CreateDatabase(mClient)
		if err != nil {
			fmt.Println("Create database error", err)
			return err
		}
	}
	return nil

}

func CreateIndex(mClient client.Client) error {
	idx, err := entity.NewIndexBinIvfFlat( // NewIndex func
		entity.HAMMING, // metricType
		1025,           // ConstructParams
	)
	if err != nil {
		return err
	}
	err = mClient.CreateIndex(
		context.Background(), // ctx
		"Images",             // CollectionName
		"descriptors",        // fieldName
		idx,                  // entity.Index
		false,                // async
	)
	if err != nil {
		return err
	}
	return nil
}

func Search(mClient client.Client, img *ImageEntity) error {
	vectorParam := img.Descriptors
	sp, err := entity.NewIndexBinIvfFlatSearchParam(5)
	if err != nil {
		return err
	}
	err = mClient.LoadCollection(context.Background(), collectionName, false)
	if err != nil {
		return err
	}
	if err != nil {
		log.Fatal("fail to create index:", err.Error())
	}

	searchResult, err := mClient.Search(
		context.Background(),     // ctx
		"Images",                 // CollectionName
		[]string{},               // partitionNames
		"",                       // expr
		[]string{"image_sha256"}, // outputFields
		[]entity.Vector{entity.BinaryVector(vectorParam)}, // vectors
		"descriptors",  // vectorField
		entity.HAMMING, // metricType
		5,              // topK
		sp,             // sp
	)
	if err != nil {
		log.Fatal("fail to search collection:", err.Error())
	}
	for _, res := range searchResult {
		fmt.Println(res.Fields.GetColumn("image_sha256"))
	}
	return nil
}
