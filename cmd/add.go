package cmd

import (
	"fmt"

	"github.com/agentx3/gofindim/data"

	"github.com/spf13/cobra"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func init() {
	addCmd.Flags().BoolVarP(&addClass, "create-class", "G", false, "Create a class")
	addCmd.Flags().StringVarP(&addDirectory, "add-directory", "A", "", "Add all images in a directory")
	RootCmd.AddCommand(addCmd)
}

var (
	addFiles     []string
	addDirectory string
	addClass     bool
	addCmd       = &cobra.Command{
		Use:   "add",
		Short: "Add files to the database",
		Long:  "Add files to the database",
		Args:  cobra.RangeArgs(0, 2),
		RunE:  ExecuteAdd,
	}
)

func ExecuteAdd(cmd *cobra.Command, args []string) error {
	// mClient, err := models.GetMilvusClient()
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	if addDirectory != "" {
		err := data.InsertDirectoryIntoWeaviate(addDirectory, client)
		if err != nil {
			return err
		}
	}

	if addClass {
		err := data.CreateImageClass(client)
		if err != nil {
			return err
		}
		return nil
	}
	for _, file := range args {
		img, err := data.NewImageFileFromPath(file)
		if err != nil {
			return err
		}
		err = data.InsertIntoWeaviate(img, client)
		if err != nil {
			return err
		}
		fmt.Printf("Added %v\n", file)
	}
	return nil

}
