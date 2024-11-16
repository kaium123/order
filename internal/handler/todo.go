package handler

import (
	"errors"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"github.com/kaium123/order/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

// TodoHandler is the request handler for the todo endpoint.
type TodoHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Find(c echo.Context) error
	FindAll(c echo.Context) error
}

type InitTodoHandler struct {
	Service service.ITodo
	Log     *log.Logger
}

type todoHandler struct {
	Handler
	service service.ITodo
	log     *log.Logger
}

// NewTodo returns a new instance of the todo handler.
func NewTodo(initTodoHandler *InitTodoHandler) TodoHandler {
	return &todoHandler{
		log:     initTodoHandler.Log,
		service: initTodoHandler.Service,
	}
}

// @Summary	Create a new todo
// @Tags		todos
// @Accept		json
// @Produce	json
// @Param		request	body		model.CreateRequest	true	"json"
// @Success	201		{object}	ResponseError{data=model.Todo}
// @Failure	400		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/todos [post]
func (t *todoHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var req model.CreateRequest
	var responseErr ResponseError

	if err := t.MustBind(c, &req); err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, err))
	}

	todo, err := t.service.Create(ctx, &req)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, err))
	}

	return c.JSON(http.StatusCreated, ResponseData{Data: todo})
}

// @Summary	Update a todo
// @Tags		todos
// @Accept		json
// @Produce	json
// @Param		body	body		model.UpdateRequestBody	true	"body"
// @Param		path	path		model.UpdateRequestPath	false	"path"
// @Success	201		{object}	ResponseData{Data=model.Todo}
// @Failure	400		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/todos/:id [put]
func (t *todoHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	var req model.UpdateRequest
	var responseErr ResponseError

	if err := t.MustBind(c, &req); err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, err))
	}

	todo, err := t.service.Update(ctx, &req)
	if err != nil {
		t.log.Error(ctx, err.Error())
		if errors.Is(err, model.ErrNotFound) {
			return c.JSON(responseErr.GetErrorResponse(http.StatusNotFound, err))
		}
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, err))
	}

	return c.JSON(http.StatusOK, ResponseData{Data: todo})
}

// @Summary	Delete a todo
// @Tags		todos
// @Param		path	path	model.DeleteRequest	false	"path"
// @Success	204
// @Failure	400	{object}	ResponseError
// @Failure	404	{object}	ResponseError
// @Failure	500	{object}	ResponseError
// @Router		/todos/:id [delete]
func (t *todoHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	var req model.DeleteRequest
	var responseErr ResponseError

	if err := t.MustBind(c, &req); err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, err))
	}

	if err := t.service.Delete(ctx, &req); err != nil {
		t.log.Error(ctx, err.Error())
		if errors.Is(err, model.ErrNotFound) {
			return c.JSON(responseErr.GetErrorResponse(http.StatusNotFound, err))
		}
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, err))
	}

	return c.JSON(http.StatusOK, "Deleted successfully")
}

// @Summary	Find a todo
// @Tags		todos
// @Param		path	path		model.FindRequest	false	"path"
// @Success	200		{object}	ResponseData{Data=model.Todo}
// @Failure	400		{object}	ResponseError
// @Failure	404		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/todos/:id [get]
func (t *todoHandler) Find(c echo.Context) error {
	ctx := c.Request().Context()
	var req model.FindRequest
	var responseErr ResponseError

	if err := t.MustBind(c, &req); err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, err))
	}

	res, err := t.service.Find(ctx, &req)
	if err != nil {
		t.log.Error(ctx, err.Error())
		if errors.Is(err, model.ErrNotFound) {
			return c.JSON(responseErr.GetErrorResponse(http.StatusNotFound, err))
		}
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, err))
	}

	return c.JSON(http.StatusOK, ResponseData{Data: res})
}

// @Summary	Find all todos
// @Tags		todos
// @Success	200	{object}	ResponseData{Data=[]model.Todo}
// @Failure	500	{object}	ResponseError
// @Router		/todos [get]
func (t *todoHandler) FindAll(c echo.Context) error {
	ctx := c.Request().Context()
	var responseErr ResponseError

	// Retrieve query parameters for 'task' and 'status'
	task := c.QueryParam("task")
	status := c.QueryParam("status")

	// Populate request params model with extracted values
	reqParams := &model.FindAllRequest{
		Task:   task,
		Status: status,
	}

	// Call the service to find all tasks based on the request params
	res, err := t.service.FindAll(ctx, reqParams)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, err))
	}

	// Return the successful result
	return c.JSON(http.StatusOK, ResponseData{Data: res})
}
