package config

import (
	"fmt"
	"github.com/kaium123/order/internal/cache"
	"github.com/kaium123/order/internal/db/bundb"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"os"
)

// Config of entire application.
type Config struct {
	Url              string        `json:"url" yaml:"url" toml:"url" mapstructure:"url"`
	DB               *bundb.Config `json:"db" yaml:"db" toml:"db" mapstructure:"db"` // nolint
	MigrateDirection Direction     `json:"migrate"`
	UI               UI
	APIServer        Server
	SwaggerServer    Server
	SQLite           SQLite
	Redis            *cache.Config
}

type LinkBoConfig struct {
	Attachments []struct {
		Content string `json:"content" yaml:"content" toml:"content" mapstructure:"content"`
		Name    string `json:"name" yaml:"name" toml:"name" mapstructure:"name"`
		Path    string `json:"path" yaml:"path" toml:"path" mapstructure:"path"`
	} `json:"attachments" yaml:"attachments" toml:"attachments" mapstructure:"attachments"`
}

type ExistingNBLSLConfig struct {
	BoPrefix     string `json:"bo_prefix" yaml:"bo_prefix" toml:"bo_prefix" mapstructure:"bo_prefix"`
	DummyImageId int64  `json:"dummy_image_id" yaml:"dummy_image_id" toml:"dummy_image_id" mapstructure:"dummy_image_id"`
}

type Oms struct {
	AppToken          string `json:"app_token" yaml:"app_token" toml:"app_token" mapstructure:"app_token"`
	BaseUrl           string `json:"base_url" yaml:"base_url" toml:"base_url" mapstructure:"base_url"`
	Transaction       string `json:"transaction" yaml:"transaction" toml:"transaction" mapstructure:"transaction"`
	ManageUserMarkets string `json:"manage_user_markets"  yaml:"manage_user_markets" toml:"manage_user_markets" mapstructure:"manage_user_markets"`
	CreateAccount     string `json:"create_account"  yaml:"create_account" toml:"create_account" mapstructure:"create_account"`
}

type IdAnalyzerConfig struct {
	ApiKey string `json:"api_key" yaml:"api_key" toml:"api_key" mapstructure:"api_key"`
	Url    string `json:"url" yaml:"url" toml:"url" mapstructure:"url"`
}

type PorichoyConfig struct {
	ApiKey              string `json:"api_key" yaml:"api_key" toml:"api_key" mapstructure:"api_key"`
	NidVerificationUrl  string `json:"nid_verification_url" yaml:"nid_verification_url" toml:"nid_verification_url" mapstructure:"nid_verification_url"`
	FaceVerificationUrl string `json:"face_verification_url" yaml:"face_verification_url" toml:"face_verification_url" mapstructure:"face_verification_url"`
	AutofillUrl         string `json:"autofill_url" yaml:"autofill_url" toml:"autofill_url" mapstructure:"autofill_url"`
}

type Services struct {
	BazarTarget     string `json:"bazar_target" yaml:"bazar_target" toml:"bazar_target" mapstructure:"bazar_target"`
	PortfolioTarget string `json:"portfolio_target" yaml:"portfolio_target" toml:"portfolio_target" mapstructure:"portfolio_target"`
	BankTarget      string `json:"bank_target" yaml:"bank_target" toml:"bank_target" mapstructure:"bank_target"`
	AuthTarget      string `json:"auth_target" yaml:"auth_target" toml:"auth_target" mapstructure:"auth_target"`
}

type Backoffice struct {
	Token          string `json:"token" yaml:"token" toml:"token" mapstructure:"token"`
	BaseUrl        string `json:"base_url" yaml:"base_url" toml:"base_url" mapstructure:"base_url"`
	Accounts       string `json:"accounts" yaml:"accounts" toml:"accounts" mapstructure:"accounts"`
	AccountDetails string `json:"account_details" yaml:"account_details" toml:"account_details" mapstructure:"account_details"`

	AccountCodePrefix string `json:"account_code_prefix" yaml:"account_code_prefix" toml:"account_code_prefix" mapstructure:"account_code_prefix"`
	DpId              string `json:"dp_id" yaml:"dp_id" toml:"dp_id" mapstructure:"dp_id"`
	AssignBo          string `json:"assign_bo" yaml:"assign_bo" toml:"assign_bo" mapstructure:"assign_bo"`
	RoutingNoDBBL     string `json:"routing_no_dbbl" yaml:"routing_no_dbbl" toml:"routing_no_dbbl" mapstructure:"routing_no_dbbl"`
	RoutingNoNBL      string `json:"routing_no_nbl" yaml:"routing_no_nbl" toml:"routing_no_nbl" mapstructure:"routing_no_nbl"`
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

// New default configurations.
func New() (conf *Config) {
	conf = new(Config)
	return
}

type Direction string

const (
	DirectionUp   Direction = "up"
	DirectionDown Direction = "down"
)

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
	c.MigrateDirection = Direction(migrate)
	if err = c.MigrateDirection.Check(); err != nil {
		panic(err)
	}

	return c

}

// MigrationDirectionFlag returns migration direction and migrateOnly flag
func (c *Config) MigrationDirectionFlag() (
	migrateDirection Direction, migrateOnly bool) {
	if c.MigrateDirection == "" {
		return DirectionUp, false
	}

	return c.MigrateDirection, true
}

// Check checks Direction values
func (m Direction) Check() (err error) {
	if m != DirectionUp && m != DirectionDown && m != "" {
		return fmt.Errorf("migration flag is not up or down: %s", m)
	}

	return
}
