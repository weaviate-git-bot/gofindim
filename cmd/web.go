package cmd

import (
	"github.com/agentx3/gofindim/web"

	"github.com/spf13/cobra"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func init() {
	RootCmd.AddCommand(webCmd)
}

var (
	webCmd = &cobra.Command{
		Use:   "web",
		Short: "web files in the database",
		Long:  "If passed a valid path to an image, it will attempt to use it to web for similar images. Otherwise, it will use a text search.",
		Args:  cobra.NoArgs,
		RunE:  Executeweb,
	}
)

func Executeweb(cmd *cobra.Command, args []string) error {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return err
	}
	web.Start(client)
	return nil
}
