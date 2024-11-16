package config

import (
	"github.com/kaium123/order/internal/cache"
	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/db/bundb"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"os"
)

// Config of entire application.
type Config struct {
	Url              string        `json:"url" yaml:"url" toml:"url" mapstructure:"url"`
	DB               *bundb.Config `json:"db" yaml:"db" toml:"db" mapstructure:"db"` // nolint
	MigrateDirection db.Direction  `json:"migrate"`
	APIServer        Server        `json:"api_server" yaml:"api_server" toml:"api_server" mapstructure:"api_server"`
	SwaggerServer    Server        `json:"swagger_server" yaml:"swagger_server" toml:"swagger_server" mapstructure:"swagger_server"`
	Redis            *cache.Config `json:"redis" yaml:"redis" toml:"redis" mapstructure:"redis"`
}

// Server is the configuration for the server.
type Server struct {
	Enable bool `json:"enable" yaml:"enable" toml:"enable" mapstructure:"enable"`
	Port   int  `json:"port" yaml:"port" toml:"port" mapstructure:"port"`
}

// New default configurations.
func New() (conf *Config) {
	conf = new(Config)
	return
}

// Load the Config from configuration files. This method panics on error.
func (c *Config) Load() *Config {

	consulPath := os.Getenv("ME_CONSUL_PATH")
	consulURL := os.Getenv("ME_CONSUL_URL")

	viper.AddRemoteProvider("consul", consulURL, consulPath)
	viper.SetConfigType("yaml") // Need to explicitly set this to json

	err := viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	viper.Unmarshal(c)

	migrate := c.MigrateDirection
	c.MigrateDirection = db.Direction(migrate)
	if err = c.MigrateDirection.Check(); err != nil {
		panic(err)
	}

	return c

}

// MigrationDirectionFlag returns migration direction and migrateOnly flag
func (c *Config) MigrationDirectionFlag() (
	migrateDirection db.Direction, migrateOnly bool) {
	if c.MigrateDirection == "" {
		return db.DirectionUp, false
	}

	return c.MigrateDirection, true
}
