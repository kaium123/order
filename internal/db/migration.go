package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/kaium123/order/internal/config/sqlxdb"
	"io/fs"
	"log"
	"net/http"
)

type Direction string

const (
	DirectionUp   Direction = "up"
	DirectionDown Direction = "down"
)

// Check checks Direction values
func (m Direction) Check() (err error) {
	if m != DirectionUp && m != DirectionDown && m != "" {
		return fmt.Errorf("migration flag is not up or down: %s", m)
	}

	return
}

// SQLFromUrl creates sql.DB instance from url
func SQLFromUrl(url string) (*sql.DB, error) {
	cfg := &sqlxdb.Config{URL: url}
	db, err := sqlxdb.New(cfg)
	if err != nil {
		return nil, err
	}

	return db.DB.DB, nil
}

// MigrateFromFS performs migration from fs.FS source, wraps MigrateFromSource
func MigrateFromFS(db *sql.DB, direction Direction, database string,
	files fs.FS) (err error) {
	src, err := httpfs.New(http.FS(files), "migrations")
	if err != nil {
		return err
	}

	return MigrateFromSource(db, direction, database, "httpfs", src)
}

// MigrateFromSource performs migration from source.Driver
func MigrateFromSource(db *sql.DB, direction Direction, database string,
	source string, files source.Driver) (err error) {

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create Postgres instance: %v", err)
	}

	m, err := migrate.NewWithInstance(source, files, database, driver)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	version, _, err := m.Version()
	if err != nil {
		fmt.Println("error ", err)
		return err
	}

	fmt.Println("current migration version -- ", version)

	log.Println("Running migration...")
	switch direction {
	case DirectionUp:
		err = m.Up()
	case DirectionDown:
		err = m.Down()
	}

	// don't emit errors on no changes made
	if err == migrate.ErrNoChange || err == migrate.ErrNilVersion {
		log.Println("No changes were made during the migration")
		return nil
	}

	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Migration applied successfully.")
	return nil
}
