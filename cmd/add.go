package cmd

import (
	models "_x3/sqldb/data"
	"_x3/sqldb/utils"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

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
		images := []*models.ImageFile{}
		filepath.WalkDir(addDirectory, func(path string, d fs.DirEntry, err error) error {
			if strings.HasPrefix(d.Name(), ".") && d.IsDir() {
				fmt.Println("Skipping directory", path)
				return filepath.SkipDir
			}
			if utils.IsImage(path) {
				img, err := models.NewImageFileFromPath(path)
				if err != nil {
					return fmt.Errorf("error creating image file from path: %w", err)
				}
				images = append(images, img)
			}
			return nil
		})
		err := models.InsertMultipleIntoWeaviate(images, client)
		if err != nil {
			return fmt.Errorf("error inserting multiple into weaviate: %w", err)
		}
		fmt.Println("Inserted multiple into weaviate")
		return nil
	}
	if addClass {
		err := models.CreateImageClass(client)
		if err != nil {
			return err
		}
		return nil
	}
	for _, file := range args {
		img, err := models.NewImageFileFromPath(file)
		if err != nil {
			return err
		}
		err = models.InsertIntoWeaviate(img, client)
		if err != nil {
			return err
		}
		fmt.Printf("Added %v\n", file)
	}
	return nil

}
