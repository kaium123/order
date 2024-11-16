// Package cmd provides the command line interface for the application.
package main

import (
	"github.com/kaium123/order/internal/config"
	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/log"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
func migrate() *cobra.Command {

	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
		Run: func(_ *cobra.Command, _ []string) {
			conf := config.New().Load()
			logger := log.New()
			// database
			postgres, err := db.New(conf.DB, logger)
			if err != nil {
				panic(err)
			}
			defer func(postgres *db.DB) {
				_ = postgres.DB.Close()
			}(postgres)
		},
	}
	return migrateCmd
}
