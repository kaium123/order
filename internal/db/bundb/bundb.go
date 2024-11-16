package bundb

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

// Config for DB.
type Config struct {
	URL                   string        `json:"url" yaml:"ulr" toml:"url" mapstructure:"url"`                                                                                 // nolint
	DialTimeout           time.Duration `json:"dial_timeout" yaml:"dial_timeout" toml:"dial_timeout" mapstructure:"dial_timeout"`                                             // nolint
	IdleTimeout           time.Duration `json:"idle_timeout" yaml:"idle_timeout" toml:"idle_timeout" mapstructure:"idle_timeout"`                                             // nolint
	ReadTimeout           time.Duration `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout" mapstructure:"read_timeout"`                                             // nolint
	WriteTimeout          time.Duration `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout" mapstructure:"write_timeout"`                                         // nolint
	RetryStatementTimeout bool          `json:"retry_statement_timeout" yaml:"retry_statement_timeout" toml:"retry_statement_timeout" mapstructure:"retry_statement_timeout"` // nolint
	MaxRetries            int           `json:"max_retries" yaml:"max_retries" toml:"max_retries" mapstructure:"max_retries"`                                                 // nolint
	MaxRetryBackoff       time.Duration `json:"max_retry_backoff" yaml:"max_retry_backoff" toml:"max_retry_backoff" mapstructure:"max_retry_backoff"`                         // nolint
	PoolSize              int           `json:"pool_size" yaml:"pool_size" toml:"pool_size" mapstructure:"pool_size"`                                                         // nolint
	PoolTimeout           time.Duration `json:"pool_timeout" yaml:"pool_timeout" toml:"pool_timeout" mapstructure:"pool_timeout"`                                             // nolint
}

// NewConfig returns new default configurations.
func NewConfig() (conf *Config) {
	conf = new(Config)
	// no default DB URL
	// all other are zero by default: i.e. use defaults
	return
}

// Options for underlying DB ORM by these configurations.
func (c *Config) options() (opts *pgdriver.Option, err error) {
	// cfg := pgdriver.Config{
	// 	Network:     "tcp",
	// 	Addr:        c.URL,
	// 	TLSConfig:   &tls.Config{InsecureSkipVerify: true},
	// 	DialTimeout: c.DialTimeout,
	// 	ReadTimeout: c.ReadTimeout,
	// }
	//
	// if opts, err = pg.ParseURL(c.URL); err != nil {
	// 	return
	// }
	// opts.DialTimeout = c.DialTimeout
	// opts.IdleTimeout = c.IdleTimeout
	// opts.ReadTimeout = c.ReadTimeout
	// opts.WriteTimeout = c.WriteTimeout
	// opts.RetryStatementTimeout = c.RetryStatementTimeout
	// opts.MaxRetries = c.MaxRetries
	// opts.MaxRetryBackoff = c.MaxRetryBackoff
	// opts.PoolSize = c.PoolSize
	// opts.PoolTimeout = c.PoolTimeout

	return
}

// The DB represents DB.
type DB struct {
	*bun.DB // underlying go-pg DB instance, embed
}

// New DB with given configurations and log.
func New(conf *Config) (db *DB, err error) {
	// Create the connection with the provided configurations.
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		pgdriver.WithDSN(conf.URL),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),
	)

	// Create the underlying SQL database connection.
	sqlDB := sql.OpenDB(pgconn)

	// Ping the database to ensure the connection is successful.
	if err = sqlDB.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool limits
	sqlDB.SetMaxIdleConns(conf.PoolSize)
	sqlDB.SetMaxOpenConns(conf.PoolSize)
	sqlDB.SetConnMaxLifetime(conf.IdleTimeout)

	// Initialize the Bun DB instance.
	db = &DB{}
	db.DB = bun.NewDB(sqlDB, pgdialect.New())

	// Verify Bun DB connection by pinging.
	if err = db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping bun DB: %w", err)
	}

	// Add query hooks for logging queries.
	db.DB.AddQueryHook(db)
	return
}

// Ping DB server.
func (db *DB) Ping(ctx context.Context) error {
	return db.DB.Ping()
}

// BeforeQuery hook (no-op for now).
func (db *DB) BeforeQuery(ctx context.Context, qe *bun.QueryEvent) context.Context {
	// Add any pre-query actions here if needed.
	return ctx
}

// AfterQuery hook (log the query operation and any errors).
func (db *DB) AfterQuery(ctx context.Context, qe *bun.QueryEvent) {
	fmt.Println(qe.Query)
}