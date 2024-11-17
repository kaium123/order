// Package server provides the API server for the application.
package server

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kaium123/order/internal/cache"
	"github.com/kaium123/order/internal/common"
	"github.com/kaium123/order/internal/config"
	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/db/bundb"
	"github.com/kaium123/order/internal/handler"
	"github.com/kaium123/order/internal/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"sync"
)

// orderAPIServer is the API server for Order
type orderAPIServer struct {
	port   int
	engine *echo.Echo
	log    *log.Logger
	db     *db.DB
}

// OrderAPIServerOpts is the options for the OrderAPIServer
type OrderAPIServerOpts struct {
	ListenPort int
	Config     config.Config
}

type InitNewAPI struct {
	OrderAPIServerOpts OrderAPIServerOpts
	Log                *log.Logger
}

var (
	dbInstance  *db.DB
	dbOnce      sync.Once
	redisClient *redis.Client
	redisOnce   sync.Once
)

// getDatabaseInstance ensures a singleton database instance
func getDatabaseInstance(ctx context.Context, config *bundb.Config, logger *log.Logger) (*db.DB, error) {
	var err error
	dbOnce.Do(func() {
		dbInstance, err = db.New(config, logger)
		if err != nil {
			logger.Error(ctx, "failed to initialize database", zap.Error(err))
			return
		}
		// Perform a ping to verify the database connection
		if pingErr := dbInstance.Ping(); pingErr != nil {
			logger.Error(ctx, fmt.Sprintf("Failed to ping the database: %v", pingErr))
			err = pingErr
			return
		}
		logger.Info(ctx, "Database ping successful")
	})
	return dbInstance, err
}

// getRedisClientInstance ensures a singleton Redis client
func getRedisClientInstance(config *cache.Config) *redis.Client {
	redisOnce.Do(func() {
		redisClient = cache.New(config)
	})
	return redisClient
}

// NewAPI initializes the API server with singleton database and Redis instances
func NewAPI(ctx context.Context, init *InitNewAPI) (Server, error) {
	// Singleton database instance
	dbInstance, err := getDatabaseInstance(ctx, init.OrderAPIServerOpts.Config.DB, init.Log)
	if err != nil {
		return nil, err
	}

	// Singleton Redis client instance
	redisClient := getRedisClientInstance(init.OrderAPIServerOpts.Config.Redis)

	// Initialize Echo server
	engine := echo.New()
	engine.HideBanner = true
	engine.HidePort = true

	// Register handlers
	handler.Register(&handler.ServiceRegistry{
		EchoEngine:  engine,
		DBInstance:  dbInstance,
		RedisClient: redisClient,
		Log:         init.Log,
	})

	// Add middleware
	engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	engine.Use(requestLogger())

	// Create API server instance
	s := &orderAPIServer{
		port:   init.OrderAPIServerOpts.ListenPort,
		engine: engine,
		log:    init.Log,
		db:     dbInstance,
	}

	// Ensure database connection is closed when the server shuts down
	go func() {
		<-ctx.Done()
		dbInstance.DB.Close()
	}()

	return s, nil
}

func (s *orderAPIServer) Name() string {
	return "orderAPIServer"
}

// Run starts the Order API server
func (s *orderAPIServer) Run() error {
	s.log.Info(context.Background(), fmt.Sprintf("%s %s serving on port %d", s.Name(), common.GetVersion(), s.port))
	return s.engine.Start(fmt.Sprintf(":%d", s.port))
}

// Shutdown stops the Order API server
func (s *orderAPIServer) Shutdown(ctx context.Context) error {
	s.log.Info(context.Background(), fmt.Sprintf("shuting down %s %s serving on port %d", s.Name(), common.GetVersion(), s.port))
	return s.engine.Shutdown(ctx)
}
