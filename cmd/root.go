package main

import (
	"github.com/spf13/cobra"
	l "log"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(serve())

	if err := rootCmd.Execute(); err != nil {
		l.Println(err)
		os.Exit(1)
	}
}
