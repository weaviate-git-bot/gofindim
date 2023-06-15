package cmd

import (
	"_x3/sqldb/ai"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(tokenizeCmd)
}

var (
	tokenizeStr []string
	tokenizeCmd = &cobra.Command{
		Use:   "tokenize",
		Short: "Tokenize a string",
		Long:  "Tokenize a string",
		Args:  cobra.ArbitraryArgs,
		RunE:  ExecuteTokenize,
	}
)

func ExecuteTokenize(cmd *cobra.Command, args []string) error {
	combined := strings.Join(args, " ")
	tokens, err := ai.Tokenize(combined)
	if err != nil {
		return err
	}
	for _, token := range tokens {
		println(token)
	}
	return nil
}
