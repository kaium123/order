// Package repository provides the database operations for the todo endpoint.
package repository

import (
	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/log"
)

// ITodo Todo is the repository for the todo endpoint.
type ITodo interface {
}

type InitTodoRepository struct {
	Db  *db.DB
	Log *log.Logger
}

type todoReceiver struct {
	log *log.Logger
	db  *db.DB
}

// NewTodo returns a new instance of the todo repository.
func NewTodo(initTodoRepository *InitTodoRepository) ITodo {
	return &todoReceiver{
		log: initTodoRepository.Log,
		db:  initTodoRepository.Db,
	}
}
