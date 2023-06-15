package cmd

import (
	"_x3/sqldb/ai"
	"fmt"
	"image"
	"os"

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
	imgFile, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}
	vector, err := ai.MakeImageEmbedding(img)
	if err != nil {
		return err
	}
	fmt.Printf("Vector: %v\n", vector)
	return nil
}
