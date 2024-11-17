package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"github.com/kaium123/order/internal/service"
	"github.com/kaium123/order/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

// AuthHandler is the request handler for the Auth endpoint.
type AuthHandler interface {
	Login(c echo.Context) error
	Logout(c echo.Context) error
}

type InitAuthHandler struct {
	Service service.IAuth
	Log     *log.Logger
}

type authHandler struct {
	Handler
	service service.IAuth
	log     *log.Logger
}

// NewAuth returns a new instance of the Auth handler.
func NewAuth(initAuthHandler *InitAuthHandler) AuthHandler {
	return &authHandler{
		log:     initAuthHandler.Log,
		service: initAuthHandler.Service,
	}
}

// Login method to authenticate the user
func (t *authHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req model.UserLoginRequest
	var responseErr utils.ResponseError

	// Bind the login request data
	if err := t.MustBind(c, &req); err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, map[string][]string{"invalid_request": []string{err.Error()}}, "Please provide a valid request body"))
	}

	// Call the service to handle login
	token, err := t.service.Login(ctx, &req)
	fmt.Println(token)
	fmt.Println(err)
	if err != nil {
		t.log.Error(ctx, err.Error())
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, nil, "The user credentials were incorrect."))
		}
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, map[string][]string{"invalid_request": []string{err.Error()}}, "The user credentials were incorrect."))
	}

	// Return the JWT or authentication token in utils
	return c.JSON(http.StatusOK, token)
}

// Logout method to invalidate the user session
func (t *authHandler) Logout(c echo.Context) error {
	ctx := c.Request().Context()
	var responseErr utils.ResponseError
	userId, err := GetUserId(c)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, nil, "Unauthorized"))
	}

	// Call the service to handle logout
	err = t.service.Logout(ctx, userId)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, map[string][]string{"invalid_request": []string{err.Error()}}, ""))
	}

	// Return success message for logout
	return c.JSON(http.StatusOK, "Logged out successfully")
}
