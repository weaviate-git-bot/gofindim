package cmd

import (
	"_x3/sqldb/models"
	"fmt"

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
	// mClient, err := models.GetMilvusClient()
	img1, err := models.NewImageFileFromPath(args[0])
	if err != nil {
		return err
	}
	img2, err := models.NewImageFileFromPath(args[1])
	if err != nil {
		return err
	}
	distance, err := models.CompareImageOrb(img1, img2)
	if err != nil {
		return err
	}
	fmt.Printf("Distance between %v and %v is %v\n", img1.Name, img2.Name, distance)
	return nil

}
