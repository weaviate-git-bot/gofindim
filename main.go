package main

import (
	"_x3/sqldb/cmd"
	"fmt"
	"os"
)

func main() {
	// Execute the Cobra command
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
