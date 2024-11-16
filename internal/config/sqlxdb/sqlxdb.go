// Package sqlxdb wraps sqlx DB ORM.
package sqlxdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	_ "github.com/jackc/pgx/v4/stdlib" // driver
	"github.com/jmoiron/sqlx"
)

const DriverName = "pgx"

// Config of the db
type Config struct {
	URL string `mapstructure:"url" yaml:"url"`
}

// NewConfig with defaults.
func NewConfig() (conf *Config) {
	conf = new(Config)
	conf.URL = "postgresql://root@cockroach:26257?sslmode=disable"
	return
}

// DB is *sqlx.DB wrapper.
type DB struct {
	*sqlx.DB
}

// New DB by given configurations.
func New(conf *Config) (db *DB, err error) {

	var conn *sqlx.DB
	if conn, err = sqlx.Open(DriverName, conf.URL); err != nil {
		return nil, fmt.Errorf("failed to db open: %w", err)
	}
	db = &DB{
		DB: conn,
	}

	if err = db.HealthCheck(context.TODO()); err != nil {
		_ = db.Close() // ignore error
		return nil, fmt.Errorf("failed to health check: %w", err)
	}

	return
}

// The HealthCheck pings the DB server.
func (db *DB) HealthCheck(ctx context.Context) error {
	return db.PingContext(ctx)
}

// TxFunc represents transaction function.
type TxFunc func(ctx context.Context, tx *sqlx.Tx) (err error)

// InTx preforms given function passing a read-write transaction to it.
func (db *DB) InTx(ctx context.Context, txFunc TxFunc) (
	err error) {

	var tx *sqlx.Tx
	if tx, err = db.BeginTxx(ctx, nil); err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback() //nolint
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(ctx, tx)
	return
}

// InReadOnlyTx preforms given function passing a read-only transaction to it.
func (db *DB) InReadOnlyTx(ctx context.Context, txFunc TxFunc) (
	err error) {

	var tx *sqlx.Tx
	tx, err = db.BeginTxx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback() //nolint
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(ctx, tx)
	return
}

// Close the DB.
func (db *DB) Close() error {
	return db.DB.Close()
}

// NotFoundError returns given notFoundError if the given err is
// sql.ErrNoRows. Otherwise, it return the err. Use it in queries
// with corresponding 'not found error' from your model.
//
//     if user, err = mydb.User(ctx, query); err != nil {
//         return nil, db.NotFoundError(err, model.ErrUserNotFound)
//     }
//     return
//
// Or any other case.
//
// Don't propagate the database/slq to a controller.
func NotFoundError(err, notFoundError error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return notFoundError
	}
	return err
}

// IsUniqueViolation returns true if given error represents unique
// constraint violation with given constrain name.
func IsUniqueViolation(err error, name string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) &&
		pgErr.Code == pgerrcode.UniqueViolation &&
		pgErr.ConstraintName == name
}

// DuplicateError returns the duplicateError if given error is unique violation
// for given constraint.
func DuplicateError(err error, constraint string, duplicateError error) error {
	if err == nil {
		return nil
	}

	if IsUniqueViolation(err, constraint) {
		return duplicateError
	}

	return err
}

// compiler time check
var (
	_ DBer = (*DB)(nil)      // must satisfy
	_ DBer = (*sqlx.DB)(nil) // must satisfy
	_ DBer = (*sqlx.Tx)(nil) // must satisfy
)

// DBer represents interface that satisfies both *sqlx.DB and *sqlx.Tx.
// Other words it's set of methods of the *Tx and the *DB intersection.
type DBer interface {
	sqlx.Queryer
	sqlx.QueryerContext
	sqlx.Preparer
	sqlx.PreparerContext

	sqlx.Execer
	sqlx.ExecerContext

	BindNamed(query string, arg interface{}) (string, []interface{}, error)

	DriverName() string

	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string,
		args ...interface{}) error

	MustExec(query string, args ...interface{}) sql.Result
	MustExecContext(ctx context.Context, query string,
		args ...interface{}) sql.Result

	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (
		sql.Result, error)

	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)

	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt,
		error)

	Rebind(query string) string

	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string,
		args ...interface{}) error
}
