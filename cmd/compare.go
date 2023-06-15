package cmd

import (
	"_x3/sqldb/data"
	"_x3/sqldb/math"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(compareCmd)
}

var (
	compareFiles []string
	compareCmd   = &cobra.Command{
		Use:   "compare",
		Short: "Compare two files",
		Long:  "Compare two files",
		Args:  cobra.ExactArgs(2),
		RunE:  ExecuteCompare,
	}
)

func ExecuteCompare(cmd *cobra.Command, args []string) error {
	// mClient, err := data.GetMilvusClient()
	_, err0 := os.Stat(args[0])
	_, err1 := os.Stat(args[1])
	if err0 == nil {
		if err1 == nil {

			img1, err := data.NewImageFileFromPath(args[0])
			if err != nil {
				return err
			}
			img2, err := data.NewImageFileFromPath(args[1])
			if err != nil {
				return err
			}
			images := []*data.ImageFile{img1, img2}
			vectors, err := data.VectorizeImages(images)
			if err != nil {
				return err
			}
			similarities := math.CosineSimilarity(vectors[0], vectors[1])
			fmt.Printf("Similarity between %v and %v is %v\n", args[0], args[1], similarities)
			return nil
		}
	}
	imageFile, err := data.NewImageFileFromPath(args[0])
	if err != nil {
		return err
	}
	textVec, imgVec, err := data.VectorizeTextAndImage(strings.Join(args[1:], " "), imageFile)
	similarities := math.CosineSimilarity(textVec, imgVec)
	fmt.Printf("Similarity between text and image is %v\n", similarities)

	return nil

}
