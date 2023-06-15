package cmd

import (
	"_x3/sqldb/models"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(dbCmd)
}

var (
	dbCreateCollection string
	dbCreateDatabase   string
	dbCreateIndex      bool
	dbDropCollection   string
	dbListCollections  bool
	dbListDatabases    bool
	dbName             string
	dbCmd              = &cobra.Command{
		Use:   "database",
		Short: "Database operations",
		Long:  "Database operations",
		Args:  cobra.RangeArgs(0, 1),
		RunE:  ExecuteDatabase,
	}
)

func ExecuteDatabase(cmd *cobra.Command, args []string) error {
	var err error
	mClient, err := models.GetMilvusClient()
	if err != nil {
		return err
	}
	if dbCreateDatabase != "" {
		err = models.CreateDatabase(mClient)
		if err != nil {
			return err
		}
		fmt.Println("Database created")
		return nil
	}
	if dbName == "" {
		return fmt.Errorf("Database name is required")
	}
	err = mClient.UsingDatabase(context.Background(), dbName)
	if err != nil {
		return err
	}
	switch {
	case dbDropCollection != "":
		err = mClient.DropCollection(context.Background(), "Images")
		if err != nil {
			return err
		}
		return nil
	case dbCreateIndex:
		err = models.CreateIndex(mClient)
		if err != nil {
			return err
		}
		fmt.Println("Index created")
		return nil

	case dbListDatabases:
		res, err := mClient.ListDatabases(context.Background())
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	case dbCreateCollection != "":
		err = models.CreateCollection(mClient)
		if err != nil {
			return err
		}
		return nil
	case dbListCollections:
		err = models.ViewCollections(mClient)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
