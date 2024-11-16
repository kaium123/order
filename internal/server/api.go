// Package server provides the API server for the application.
package server

import (
	"context"
	"fmt"
	"github.com/kaium123/order/internal/cache"
	"github.com/kaium123/order/internal/common"
	"github.com/kaium123/order/internal/config"
	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/handler"
	"github.com/kaium123/order/internal/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func NewAPI(ctx context.Context, init *InitNewAPI) (Server, error) {
	// database
	dbInstance, err := db.New(init.OrderAPIServerOpts.Config.DB, init.Log)
	if err != nil {
		panic(err)
	}

	// Perform a ping to verify the database connection is working
	err = dbInstance.Ping()
	if err != nil {
		init.Log.Error(ctx, fmt.Sprintf("Failed to ping the database: %v", err))
		return nil, err
	} else {
		init.Log.Info(ctx, "Database ping successful")
	}

	// Initialize other components like Redis and Echo server
	redisClient := cache.New(init.OrderAPIServerOpts.Config.Redis)

	engine := echo.New()
	engine.HideBanner = true
	engine.HidePort = true

	handler.Register(&handler.ServiceRegistry{
		EchoEngine:  engine,
		DBInstance:  dbInstance,
		RedisClient: redisClient,
		Log:         init.Log,
	})

	engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	engine.Use(requestLogger())

	s := &orderAPIServer{
		port:   init.OrderAPIServerOpts.ListenPort,
		engine: engine,
		log:    init.Log,
		db:     dbInstance,
	}

	// Closing the database connection when the server shuts down
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
