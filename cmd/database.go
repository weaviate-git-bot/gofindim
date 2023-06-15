package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func init() {
	dbCmd.Flags().StringVarP(&dbName, "", "n", "", "Database name")
	dbCmd.Flags().BoolVarP(&dbListCollections, "list-collections", "l", false, "List collections")
	dbCmd.Flags().BoolVarP(&dbListDatabases, "list-databases", "L", false, "List databases")
	dbCmd.Flags().BoolVarP(&dbCreateIndex, "create-index", "I", false, "Create index")
	dbCmd.Flags().StringVarP(&dbDropCollection, "drop-collection", "D", "", "Drop a collection")
	dbCmd.Flags().StringVarP(&dbCreateCollection, "create-collection", "c", "", "Create a collection")
	dbCmd.Flags().StringVarP(&dbCreateDatabase, "create-database", "C", "", "Create a database")
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
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	if dbDropCollection != "" {

		client.Schema().ClassDeleter().WithClassName(dbDropCollection).Do(context.Background())
		fmt.Printf("Dropped collection %v\n", dbDropCollection)

	}
	return nil
}
