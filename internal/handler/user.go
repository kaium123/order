package handler

import (
	"errors"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"github.com/kaium123/order/internal/service"
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
	var responseErr ResponseError

	// Bind the login request data
	if err := t.MustBind(c, &req); err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, err))
	}

	// Call the service to handle login
	token, err := t.service.Login(ctx, &req)
	if err != nil {
		t.log.Error(ctx, err.Error())
		if errors.Is(err, model.ErrInvalidCredentials) {
			return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, err))
		}
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, err))
	}

	// Return the JWT or authentication token in response
	return c.JSON(http.StatusOK, ResponseData{Data: token})
}

// Logout method to invalidate the user session
func (t *authHandler) Logout(c echo.Context) error {
	ctx := c.Request().Context()
	var responseErr ResponseError
	userId, err := GetUserId(c)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, err))
	}

	// Call the service to handle logout
	err = t.service.Logout(ctx, userId)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, err))
	}

	// Return success message for logout
	return c.JSON(http.StatusOK, "Logged out successfully")
}
