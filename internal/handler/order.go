package handler

import (
	"errors"
	"fmt"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"github.com/kaium123/order/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// OrderHandler is the request handler for the Order endpoint.
type OrderHandler interface {
	CreateOrder(c echo.Context) error
	CancelOrder(c echo.Context) error
	FindAllOrders(c echo.Context) error
}

type InitOrderHandler struct {
	Service service.IOrder
	Log     *log.Logger
}

type orderHandler struct {
	Handler
	service service.IOrder
	log     *log.Logger
}

// NewOrder returns a new instance of the Order handler.
func NewOrder(initOrderHandler *InitOrderHandler) OrderHandler {
	return &orderHandler{
		log:     initOrderHandler.Log,
		service: initOrderHandler.Service,
	}
}

func (t *orderHandler) CreateOrder(c echo.Context) error {
	ctx := c.Request().Context()
	var responseErr ResponseError
	var req model.Order

	// Retrieve the user_id from the context
	userId, err := GetUserId(c)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, err))
	}

	if err := t.MustBind(c, &req); err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, err))
	}

	req.UserID = userId
	Order, err := t.service.CreateOrder(ctx, &req)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, err))
	}

	return c.JSON(http.StatusCreated, ResponseData{Data: Order})
}

func (t *orderHandler) CancelOrder(c echo.Context) error {
	ctx := c.Request().Context()
	var req model.OrderCancelRequest
	var responseErr ResponseError
	// Retrieve the user_id from the context
	userId, err := GetUserId(c)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, err))
	}

	// Get consignmentId directly from query parameter
	consignmentId := c.Param("CONSIGNMENT_ID")
	fmt.Println(consignmentId)
	req.ConsignmentID = consignmentId
	req.UserId = userId

	if err := t.MustBind(c, &req); err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, err))
	}

	if err := t.service.CancelOrder(ctx, &req); err != nil {
		t.log.Error(ctx, err.Error())
		if errors.Is(err, model.ErrNotFound) {
			return c.JSON(responseErr.GetErrorResponse(http.StatusNotFound, err))
		}
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, err))
	}

	return c.JSON(http.StatusOK, "Deleted successfully")
}

func (t *orderHandler) FindAllOrders(c echo.Context) error {
	ctx := c.Request().Context()
	var responseErr ResponseError
	userId, err := GetUserId(c)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, err))
	}

	// Retrieve query parameters for 'task' and 'status'
	transferStatus := c.QueryParam("transfer_status")
	archive, err := strconv.Atoi(c.QueryParam("archive"))
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	page, err := strconv.Atoi(c.QueryParam("page"))

	// Populate request params model with extracted values
	reqParams := &model.FindAllRequest{
		TransferStatus: transferStatus,
		Archive:        archive,
		Limit:          limit,
		Offset:         (page - 1) * limit,
		UserId:         userId,
	}

	// Call the service to find all tasks based on the request params
	res, err := t.service.FindAllOrders(ctx, reqParams)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, err))
	}

	// Return the successful result
	return c.JSON(http.StatusOK, ResponseData{Data: res})
}

func GetUserId(c echo.Context) (int64, error) {

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	intUserId, err := strconv.Atoi(userID)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}
	return int64(intUserId), nil
}
