package server

import (
	"context"
	"github.com/kaium123/order/internal/log"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSwagger(t *testing.T) {
	opts := SwaggerServerOpts{
		ListenPort: 8080,
	}

	logger := log.New()
	server := NewSwagger(context.Background(), &InitNewSwagger{
		SwaggerServerOpts: opts,
		Log:               logger,
	})

	assert.NotNil(t, server)
	assert.Equal(t, "swaggerServer", server.Name())

	swaggerServer, ok := server.(*swaggerServer)
	assert.True(t, ok)
	assert.Equal(t, opts.ListenPort, swaggerServer.port)
	assert.NotNil(t, swaggerServer.engine)
	assert.NotNil(t, swaggerServer.log)

	// Check if the Echo instance has the expected routes
	routes := swaggerServer.engine.Routes()
	foundSwaggerRoute := false
	for _, route := range routes {
		if route.Path == "/swagger/*" && route.Method == echo.GET {
			foundSwaggerRoute = true
			break
		}
	}
	assert.True(t, foundSwaggerRoute, "Swagger route not found")
}
