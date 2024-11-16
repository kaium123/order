package server

import (
	"context"
	"fmt"
	"github.com/kaium123/order/internal/common"
	"github.com/kaium123/order/internal/log"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// swaggerServer is the API server for Todo
type swaggerServer struct {
	port   int
	engine *echo.Echo
	log    *log.Logger
}

// SwaggerServerOpts is the options for the swaggerServer
type SwaggerServerOpts struct {
	ListenPort int
}

type InitNewSwagger struct {
	SwaggerServerOpts SwaggerServerOpts
	Log               *log.Logger
}

// NewSwagger returns a new instance of the Swagger server
func NewSwagger(ctx context.Context, init *InitNewSwagger) Server {

	engine := echo.New()
	engine.HideBanner = true
	engine.HidePort = true

	engine.Use(requestLogger())

	engine.GET("/swagger/*", echoSwagger.WrapHandler)

	s := &swaggerServer{
		port:   init.SwaggerServerOpts.ListenPort,
		engine: engine,
		log:    init.Log,
	}

	return s
}

func (s *swaggerServer) Name() string {
	return "swaggerServer"
}

func (s *swaggerServer) Run() error {
	s.log.Info(context.Background(), fmt.Sprintf("%s %s serving on port %d", s.Name(), common.GetVersion(), s.port))
	return s.engine.Start(fmt.Sprintf(":%d", s.port))
}

func (s *swaggerServer) Shutdown(ctx context.Context) error {
	s.log.Info(context.Background(), fmt.Sprintf("shuting down %s %s serving on port %d", s.Name(), common.GetVersion(), s.port))
	return s.engine.Shutdown(ctx)
}
