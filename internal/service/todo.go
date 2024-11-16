// Package service provides the business logic for the todo endpoint.
package service

import (
	"context"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"github.com/kaium123/order/internal/repository"
)

// Todo is the service for the todo endpoint.
type ITodo interface {
	Create(ctx context.Context, reqTodo *model.CreateRequest) (*model.Todo, error)
	Update(ctx context.Context, reqTodo *model.UpdateRequest) (*model.Todo, error)
	Delete(ctx context.Context, reqParams *model.DeleteRequest) error
	Find(ctx context.Context, reqParams *model.FindRequest) (*model.Todo, error)
	FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error)
}

type todoReceiver struct {
	log            *log.Logger
	todoRepository repository.ITodo
	redisCache     repository.IRedisCache
}

func (t todoReceiver) Create(ctx context.Context, reqTodo *model.CreateRequest) (*model.Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (t todoReceiver) Update(ctx context.Context, reqTodo *model.UpdateRequest) (*model.Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (t todoReceiver) Delete(ctx context.Context, reqParams *model.DeleteRequest) error {
	//TODO implement me
	panic("implement me")
}

func (t todoReceiver) Find(ctx context.Context, reqParams *model.FindRequest) (*model.Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (t todoReceiver) FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error) {
	//TODO implement me
	panic("implement me")
}

type InitTodoService struct {
	Log            *log.Logger
	TodoRepository repository.ITodo
	RedisCache     repository.IRedisCache
}

// NewTodo creates a new Todo service.
func NewTodo(initTodoService *InitTodoService) ITodo {
	return &todoReceiver{
		log:            initTodoService.Log,
		todoRepository: initTodoService.TodoRepository,
		redisCache:     initTodoService.RedisCache,
	}
}
