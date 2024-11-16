package model

import (
	"errors"
	"fmt"
	"time"
)

// Todo is the model for the todo endpoint.
type Todo struct {
	ID        int `gorm:"primaryKey"`
	Task      string
	Status    Status
	Priority  TodoPriority
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type FindAllRequest struct {
	Task   string
	Status string
}

// UpdateRequestPath is the request parameter for updating a todo
type UpdateRequestPath struct {
	ID int `param:"id" validate:"required"`
}

// UpdateRequest is the request parameter for updating a todo
type UpdateRequest struct {
	UpdateRequestBody
	UpdateRequestPath
}

// CreateRequest is the request parameter for creating a new todo
type CreateRequest struct {
	Task     string `json:"task" validate:"required"`
	Priority string `json:"priority" validate:"required"`
}

// UpdateRequestBody is the request body for updating a todo
type UpdateRequestBody struct {
	Task   string `json:"task,omitempty"`
	Status Status `json:"status,omitempty"`
}

// DeleteRequest is the request parameter for deleting a todo
type DeleteRequest struct {
	ID int `param:"id" validate:"required"`
}

// FindRequest is the request parameter for finding a todo
type FindRequest struct {
	ID int `param:"id" validate:"required"`
}

// NewTodo returns a new instance of the todo model.
func NewTodo(req *CreateRequest) *Todo {
	return &Todo{
		Task:     req.Task,
		Status:   Created,                    // Set default status
		Priority: TodoPriority(req.Priority), // Map priority directly
	}
}

// NewUpdateTodo returns a new instance of the todo model for updating.
func NewUpdateTodo(req *UpdateRequest) *Todo {
	return &Todo{
		ID:     req.ID,
		Task:   req.Task,
		Status: req.Status,
	}
}

// Status is the status of the task.
type Status string
type TodoPriority string

const (
	// TP_Low is the priority for a task.
	TP_Low = TodoPriority("low")
	// TP_Medium is the priority for a task.
	TP_Medium = TodoPriority("medium")
	// TP_High is the priority for a task.
	TP_High = TodoPriority("high")
)

const (
	// Created is the status for a created task.
	Created = Status("created")
	// Processing is the status for a processing task.
	Processing = Status("processing")
	// Done is the status for a done task.
	Done = Status("done")
)

// StatusMap is a map of task status.
var StatusMap = map[Status]bool{
	Created:    true,
	Processing: true,
	Done:       true,
}

func (t *Todo) ValidateCreateRequest() error {
	if t.Task == "" {
		return errors.New("task cannot be empty")
	}

	if t.Priority != TP_High && t.Priority != TP_Low && t.Priority != TP_Medium {
		return errors.New("invalid priority not accepted")
	}
	// Add additional validation as needed
	return nil
}

func (t *Todo) PrepareUpdatedTodo(currentTodo *Todo) *Todo {

	// 空文字列の場合、現在の値を使用
	if t.Task == "" {
		t.Task = currentTodo.Task
	}

	fmt.Println(t.Status)

	if t.Status == "" {
		t.Status = currentTodo.Status
	}

	fmt.Println(t.Status)

	t.CreatedAt = currentTodo.CreatedAt
	t.Priority = currentTodo.Priority

	return currentTodo
}
