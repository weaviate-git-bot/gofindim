package cmd

import (
	"fmt"
	"os"

	"github.com/agentx3/gofindim/data"

	"github.com/spf13/cobra"
)

var (
	// galleryURL  string
	initialize      bool
	scanFolder      string
	imageFile       string
	findFolder      string
	similar         int
	searchFile      string
	debugOption     string
	imagesToCompare []string
	RootCmd         = &cobra.Command{
		Use:   "gofindim",
		Short: "Find similar images",
		Args:  cobra.NoArgs,
		RunE:  Execute,
	}
)

func Execute(cmd *cobra.Command, args []string) error {

	if initialize {
	}
	if len(imagesToCompare) == 2 {
		img1, err := data.NewImageFileFromPath(imagesToCompare[0])
		if err != nil {
			return err
		}
		img2, err := data.NewImageFileFromPath(imagesToCompare[1])
		if err != nil {
			return err
		}
		distance, err := data.CompareImageOrb(img1, img2)
		if err != nil {
			return err
		}
		fmt.Printf("Distance between %v and %v is %v\n", img1.Name, img2.Name, distance)
		return nil
	}
	fmt.Printf("Ending execution\n")
	return nil
}

func main() {
	// Execute the Cobra command
	RootCmd.Flags().BoolVarP(&initialize, "initialize", "I", false, "Initialize the database")
	RootCmd.Flags().StringVarP(&scanFolder, "scan", "S", "", "Scan a directory for images")
	RootCmd.Flags().StringVarP(&imageFile, "image-file", "i", "", "Image file to hash")
	RootCmd.Flags().StringVarP(&findFolder, "find-folder", "F", "", "Find similar images in a folder")
	RootCmd.Flags().StringVarP(&debugOption, "debug-opt", "D", "", "Debug option")
	RootCmd.Flags().StringArrayVarP(&imagesToCompare, "compareImages", "C", []string{}, "Compare two images")
	RootCmd.Flags().StringVarP(&searchFile, "search-file", "Z", "", "Search for a file in the database")
	RootCmd.Flags().IntVar(&similar, "similar", 0, "Image file to hash")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
