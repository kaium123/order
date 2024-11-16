// Package model provides the data models for the application.
package model

import "github.com/kaium123/order/internal/cache"

// Config is the configuration for the application.
type Config struct {
	UI            UI
	APIServer     Server
	SwaggerServer Server
	SQLite        SQLite
	Redis         *cache.Config
}

// UI is the configuration for the UI.
type UI struct {
	URL string `validate:"required"`
}

// Server is the configuration for the server.
type Server struct {
	Enable bool
	Port   int
}

// SQLite is the configuration for the SQLite database.
type SQLite struct {
	DBFilename string `validate:"required"`
}