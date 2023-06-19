package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	models "github.com/agentx3/gofindim/data"
	"github.com/agentx3/gofindim/utils"

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
		imagePaths := []string{}
		// images := []*models.ImageFile{}

		filepath.WalkDir(addDirectory, func(path string, d fs.DirEntry, err error) error {
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
			imageBatch := []*models.ImageFile{}
			if end > len(imagePaths) {
				end = len(imagePaths)
			}
			batch := imagePaths[i:end]
			for _, img := range batch {
				imageFile, err := models.NewImageFileFromPath(img)
				if err != nil {
					fmt.Printf("error creating image file from path %v: %v", img, err)
					continue
				}
				imageBatch = append(imageBatch, imageFile)
			}
			waitGroup.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				err = models.InsertMultipleIntoWeaviate(imageBatch, client)
				if err != nil {
					fmt.Printf("error inserting multiple into weaviate: %w", err)
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
