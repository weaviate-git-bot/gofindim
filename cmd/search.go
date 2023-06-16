package cmd

import (
	"_x3/sqldb/data"
	"_x3/sqldb/utils"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func init() {
	searchCmd.Flags().Float32VarP(&searchThreshold, "threshold", "t", 0.5, "Threshold for search")
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 5, "Limit for search")
	searchCmd.Flags().StringArrayVarP(&searchFields, "field", "f", []string{"name"}, "Fields to grab")
	RootCmd.AddCommand(searchCmd)
}

var (
	searchFiles     []string
	searchThreshold float32
	searchFields    []string
	searchLimit     int
	searchCmd       = &cobra.Command{
		Use:   "search",
		Short: "search files in the database",
		Long:  "If passed a valid path to an image, it will attempt to use it to search for similar images. Otherwise, it will use a text search.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  ExecuteSearch,
	}
)

func ExecuteSearch(cmd *cobra.Command, args []string) error {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return err
	}
	err = searchArgs(args, client)
	if err != nil {
		return err
	}
	return nil
}

func searchArgs(args []string, client *weaviate.Client) error {
	var err error
	var results *[]data.ImageNode
	if utils.IsImage(args[0]) {
		img, err := data.NewImageFileFromPath(args[0])
		if err != nil {
			return err
		}
		results, err = searchImage(img, client)
		if err != nil {
			return err
		}
	} else {
		searchString := strings.Join(args, " ")
		results, err = searchText(searchString, client)
		if err != nil {
			return err
		}
	}
	printResults(*results)
	return nil
}

func searchImage(img *data.ImageFile, client *weaviate.Client) (*[]data.ImageNode, error) {
	results, err := data.SearchWeaviateWithImageFile(img, searchThreshold, searchLimit, searchFields, client)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func searchText(text string, client *weaviate.Client) (*[]data.ImageNode, error) {
	results, err := data.SearchWeaviateWithText(text, searchThreshold, searchLimit, searchFields, client)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func printResults(results []data.ImageNode) error {
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return err
	}
	fmt.Println(string(resultsJSON))
	return nil
}
