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

	"go.uber.org/zap"
)

// todoAPIServer is the API server for Todo
type todoAPIServer struct {
	port   int
	engine *echo.Echo
	log    *log.Logger
	db     *db.DB
}

// TodoAPIServerOpts is the options for the TodoAPIServer
type TodoAPIServerOpts struct {
	ListenPort int
	Config     config.Config
}

type InitNewAPI struct {
	TodoAPIServerOpts TodoAPIServerOpts
	Log               *log.Logger
}

// NewAPI returns a new instance of the Todo API server
func NewAPI(ctx context.Context, init *InitNewAPI) (Server, error) {

	conf := config.New().Load()
	logger := log.New()

	fmt.Println(conf)
	// database
	dbInstance, err := db.New(conf.DB, logger)
	if err != nil {
		panic(err)
	}
	defer func(postgres *db.DB) {
		_ = postgres.DB.Close()
	}(dbInstance)

	redisClient := cache.New(init.TodoAPIServerOpts.Config.Redis)

	engine := echo.New()
	engine.HideBanner = true
	engine.HidePort = true

	handler.Register(&handler.ServiceRegistry{
		EchoEngine:  engine,
		DBInstance:  dbInstance,
		RedisClient: redisClient,
		Log:         init.Log,
	})

	allowOrigins := []string{init.TodoAPIServerOpts.Config.UI.URL}
	if init.TodoAPIServerOpts.Config.SwaggerServer.Enable {
		allowOrigins = append(allowOrigins, fmt.Sprintf("http://localhost:%d", init.TodoAPIServerOpts.Config.SwaggerServer.Port))
	}
	init.Log.Info(ctx, "CORS allowed origins: ", zap.Any("origins: ", allowOrigins))
	engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	engine.Use(requestLogger())

	s := &todoAPIServer{
		port:   init.TodoAPIServerOpts.ListenPort,
		engine: engine,
		log:    init.Log,
		db:     dbInstance,
	}
	return s, nil
}

func (s *todoAPIServer) Name() string {
	return "todoAPIServer"
}

// Run starts the Todo API server
func (s *todoAPIServer) Run() error {
	s.log.Info(context.Background(), fmt.Sprintf("%s %s serving on port %d", s.Name(), common.GetVersion(), s.port))
	return s.engine.Start(fmt.Sprintf(":%d", s.port))
}

// Shutdown stops the Todo API server
func (s *todoAPIServer) Shutdown(ctx context.Context) error {
	s.log.Info(context.Background(), fmt.Sprintf("shuting down %s %s serving on port %d", s.Name(), common.GetVersion(), s.port))
	return s.engine.Shutdown(ctx)
}
