package data

import (
	"context"
	"fmt"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func DeleteWeaviateWithUUID(ctx context.Context, client *weaviate.Client, imageUUID, path string) error {
	fmt.Println("Deleting image with UUID: ", imageUUID)

	err := client.Data().Deleter().
		WithClassName("Image").
		WithID(imageUUID).
		Do(ctx)
	if err != nil {
		fmt.Println("Error deleting image with UUID: ", imageUUID)
	} else {
		fmt.Println("Deleted image with UUID: ", imageUUID)
		err = os.Remove(fmt.Sprintf(path))
		if err != nil {
			fmt.Println("Error deleting image with UUID from disk: ", imageUUID)
		}
	}

	return err
}
