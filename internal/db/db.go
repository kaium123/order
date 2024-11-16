package db

import (
	"context"
	"database/sql"
	"github.com/kaium123/order/internal/db/bundb"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"github.com/uptrace/bun"
)

// compiler time check
var _ model.DB = (*DB)(nil)

// DB is a database representation
type DB struct {
	//DB  *pgdb.DB // underlying go-pg DB wrapper instance
	*bun.DB // underlying go-pg DB wrapper instance
	log     *log.Logger
}

// Tx represents transactions
type Tx struct {
	//orm.DB
	*bun.Tx
	log *log.Logger
}

// compiler time check
var _ model.Repository = (*Tx)(nil)

//
//// New DB with given configurations and logger.
//func New(conf *pgdb.Config, logger *log.Logger) (db *DB, err error) {
//
//	var pg *pgdb.DB
//	if pg, err = pgdb.New(conf); err != nil {
//		return
//	}
//
//	db = new(DB)
//	db.DB = pg // embed
//	db.log = logger.Named("me_pg")
//
//	db.registerTables()
//	return
//}

// New DB with given configurations and logger.
func New(conf *bundb.Config, logger *log.Logger) (db *DB, err error) {

	var pg *bundb.DB
	if pg, err = bundb.New(conf); err != nil {
		return
	}

	db = new(DB)
	db.DB = pg.DB // embed
	db.registerTables()
	db.log = logger.Named("db_model")
	db.log.Info(context.Background(), "db initialization done")
	return
}

// register all tables for relations
func (db *DB) registerTables() {

}

//
//// InTx runs given function in SQL-transaction.
//func (db *DB) InTx(ctx context.Context, txFunc model.TxFunc) (err error) {
//
//	err = db.DB.RunInTransaction(ctx, func(tx *pg.Tx) (err error) {
//		return txFunc(ctx, &Tx{tx, db.log})
//	})
//	return
//}

// InTx runs given function in SQL-transaction.
func (db *DB) InTx(ctx context.Context, txFunc model.TxFunc) (err error) {

	err = db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) (err error) {
		return txFunc(ctx, &Tx{
			&tx,
			db.log,
		})
	})
	return
}
