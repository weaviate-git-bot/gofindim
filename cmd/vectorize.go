package cmd

import (
	"github.com/agentx3/gofindim/data"
	"github.com/agentx3/gofindim/math"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(vectorizeCmd)
}

var (
	// vectorizeFile string
	vectorizeCmd = &cobra.Command{
		Use:   "vectorize",
		Short: "Vectorize an image",
		Long:  "Vectorize an image",
		Args:  cobra.ExactArgs(1),
		RunE:  ExecuteVectorize,
	}
)

func ExecuteVectorize(cmd *cobra.Command, args []string) error {
	_, err := os.Stat(args[0])
	if err == nil {
		if len(args) == 1 {
			imageFile, err := data.NewImageFileFromPath(args[0])
			if err != nil {
				return err
			}
			vector, err := data.VectorizeImage(imageFile)
			if err != nil {
				return err
			}
			fmt.Printf("%v\v", vector)
			return nil
		} else {
			imageFile, err := data.NewImageFileFromPath(args[0])
			if err != nil {
				return err
			}
			textVec, imgVec, err := data.VectorizeTextAndImage(strings.Join(args[1:], " "), imageFile)
			similarities := math.CosineSimilarity(textVec, imgVec)
			fmt.Printf("Similarity between text and image is %v\n", similarities)
			return nil
		}
	} else {
		vector, err := data.VectorizeText(strings.Join(args, " "))
		if err != nil {
			return err
		}
		fmt.Printf("%v", vector)
	}

	return nil

}
