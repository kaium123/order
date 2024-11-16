package handler

import (
	"errors"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"github.com/kaium123/order/internal/service"
	"github.com/kaium123/order/internal/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
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
	var responseErr utils.ResponseError
	var req model.Order

	// Retrieve the user_id from the context
	userId, err := GetUserId(c)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, nil, "Unauthorized"))
	}

	if err := t.MustBind(c, &req); err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, map[string][]string{"invalid_request": []string{err.Error()}}, "Please provide a valid request body"))
	}

	req.UserID = userId
	validationErr := req.Validate()
	if validationErr != nil {
		t.log.Error(ctx, "validation errors : ", zap.Any("", validationErr))
		return c.JSON(responseErr.GetErrorResponse(http.StatusUnprocessableEntity, validationErr.Errors, "Please fix the given errors"))
	}

	Order, err := t.service.CreateOrder(ctx, &req)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, map[string][]string{"order_creation_error": []string{err.Error()}}, "Internal server error"))
	}

	return c.JSON(http.StatusCreated, utils.GetResponseData(http.StatusOK, Order, "Order Created Successfully"))
}

func (t *orderHandler) CancelOrder(c echo.Context) error {
	ctx := c.Request().Context()
	var req model.OrderCancelRequest
	var responseErr utils.ResponseError

	// Retrieve the user_id from the context
	userId, err := GetUserId(c)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, nil, "Unauthorized"))
	}

	// Get consignmentId directly from query parameter
	consignmentId := c.Param("CONSIGNMENT_ID")
	req.ConsignmentID = consignmentId
	req.UserId = userId

	if err := t.MustBind(c, &req); err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusBadRequest, map[string][]string{"invalid_request": []string{err.Error()}}, "Please provide a valid request body"))
	}

	if err := t.service.CancelOrder(ctx, &req); err != nil {
		t.log.Error(ctx, err.Error())
		if errors.Is(err, model.ErrNotFound) {
			return c.JSON(responseErr.GetErrorResponse(http.StatusNotFound, map[string][]string{"order_cancellation_error": []string{err.Error()}}, "Order not found"))
		}
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, map[string][]string{"order_cancellation_error": []string{err.Error()}}, "Internal server error"))
	}

	return c.JSON(http.StatusCreated, utils.GetResponseData(http.StatusOK, nil, "Order Cancelled Successfully"))
}

func (t *orderHandler) FindAllOrders(c echo.Context) error {
	ctx := c.Request().Context()
	var responseErr utils.ResponseError

	userId, err := GetUserId(c)
	if err != nil {
		t.log.Error(ctx, err.Error())
		return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, nil, "Unauthorized"))
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
		return c.JSON(responseErr.GetErrorResponse(http.StatusInternalServerError, map[string][]string{"order_finding_error": []string{err.Error()}}, "Internal server error"))
	}

	// Return the successful result
	return c.JSON(http.StatusCreated, utils.GetResponseData(http.StatusOK, res, "Orders successfully fetched."))
}

func GetUserId(c echo.Context) (int64, error) {

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}

	return userID, nil
}
