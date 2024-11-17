package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/kaium123/order/internal/config"
	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/server"
	"github.com/kaium123/order/sql"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
)

// serve returns a new `serve` command to be used as a sub-command to root
func serve() *cobra.Command {
	serveCmd := cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Run: func(_ *cobra.Command, _ []string) {
			var (
				servers []server.Server
				conf    = config.New().Load()
				logger  = log.New()
				ctx     = context.Background()
			)

			// Run migrations
			if migrateOnly, err := runMigrations(ctx, conf, logger); migrateOnly {
				panic(err)
			}

			// Initialize API server
			initNewAPI := &server.InitNewAPI{
				OrderAPIServerOpts: server.OrderAPIServerOpts{
					ListenPort: conf.APIServer.Port,
					Config:     *conf,
				},
				Log: logger,
			}

			apiServer, err := server.NewAPI(ctx, initNewAPI)
			if err != nil {
				logger.Fatal(ctx, "failed to init API.", zap.Error(err))
				panic(err)
			}
			servers = append(servers, apiServer)

			// Initialize Swagger server if enabled
			if conf.SwaggerServer.Enable {
				initNewSwagger := &server.InitNewSwagger{
					SwaggerServerOpts: server.SwaggerServerOpts{
						ListenPort: conf.SwaggerServer.Port,
					},
					Log: logger,
				}

				swagServer := server.NewSwagger(ctx, initNewSwagger)
				servers = append(servers, swagServer)
			}

			// Handle graceful shutdown
			ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
			defer stop()

			for _, s := range servers {
				server := s
				go func() {
					if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						logger.Fatal(ctx, fmt.Sprintf("shutting down %s. ", server.Name()), zap.Error(err))
						panic(err)
					}
				}()
			}

			logger.Info(ctx, "server started")

			<-ctx.Done()
			logger.Info(ctx, "server shutting down")
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			for _, s := range servers {
				if err := s.Shutdown(ctx); err != nil {
					logger.Fatal(ctx, "error while shutting down.", zap.Error(err))
					panic(err)
				}
			}
			logger.Info(ctx, "server shutdown gracefully")
		},
	}
	return &serveCmd
}

// runMigrations handles database migrations
func runMigrations(ctx context.Context, conf *config.Config, logger *log.Logger) (bool, error) {
	migrateDirection, migrateOnly := conf.MigrationDirectionFlag()

	migrateDB, err := db.SQLFromUrl(conf.DB.URL)
	if err != nil {
		logger.Fatal(ctx, "failed to connect to database for migration", zap.Error(err))
		return false, err
	}
	defer migrateDB.Close()

	migrations := sql.GetMigrations()
	err = db.MigrateFromFS(migrateDB, migrateDirection, "orders", migrations)
	if err != nil {
		logger.Fatal(ctx, "failed to run migrations", zap.Error(err))
		return false, err
	}

	if migrateOnly {
		logger.Info(ctx, "Migration complete, exiting")
	}

	return migrateOnly, nil
}
