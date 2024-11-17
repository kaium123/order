// Package service provides the business logic for the Order endpoint.
package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"github.com/kaium123/order/internal/repository"
	"math/big"
	"time"
)

// Order is the service for the Order endpoint.
type IOrder interface {
	CreateOrder(ctx context.Context, reqOrder *model.Order) (*model.CreateOrderResponse, error)
	CancelOrder(ctx context.Context, reqParams *model.OrderCancelRequest) error
	FindAllOrders(ctx context.Context, reqParams *model.FindAllRequest) (*model.FindAllResponse, error)
}

type OrderReceiver struct {
	log             *log.Logger
	OrderRepository repository.IOrder
	redisCache      repository.IRedisCache
}

type InitOrderService struct {
	Log             *log.Logger
	OrderRepository repository.IOrder
	RedisCache      repository.IRedisCache
}

// NewOrder creates a new Order service.
func NewOrder(initOrderService *InitOrderService) IOrder {
	return &OrderReceiver{
		log:             initOrderService.Log,
		OrderRepository: initOrderService.OrderRepository,
		redisCache:      initOrderService.RedisCache,
	}
}
func (o *OrderReceiver) CreateOrder(ctx context.Context, reqOrder *model.Order) (*model.CreateOrderResponse, error) {
	// Generate consignment ID
	consignmentID := GenerateConsignmentID("DA", 6)
	reqOrder.OrderConsignmentID = consignmentID
	if reqOrder.RecipientCity == 1 {
		reqOrder.CalculateDeliveryFee(60)
	} else {
		reqOrder.CalculateDeliveryFee(100)
	}

	reqOrder.CalculateCodFee()
	reqOrder.CalculateTotalFee()
	reqOrder.OrderStatus = model.Pending
	reqOrder.DeliveryType = model.Delivery
	reqOrder.ItemType = model.Parcel
	reqOrder.OrderTypeID = 1

	// Create the order in the repository (DB)
	order, err := o.OrderRepository.CreateOrder(ctx, reqOrder)
	if err != nil {
		return nil, err
	}

	// Cache the order in Redis
	err = o.redisCache.CacheOrder(ctx, *order)
	if err != nil {
		o.log.Error(ctx, fmt.Sprintf("Failed to cache order with ID %s: %v", order.OrderConsignmentID, err))
	}

	// Return the utils
	return &model.CreateOrderResponse{
		ConsignmentID:   order.OrderConsignmentID,
		MerchantOrderID: order.MerchantOrderID,
		OrderStatus:     order.OrderStatus.String(),
		DeliveryFee:     order.DeliveryFee,
	}, nil
}

func (o *OrderReceiver) CancelOrder(ctx context.Context, reqParams *model.OrderCancelRequest) error {
	err := o.OrderRepository.CancelOrder(ctx, reqParams)
	if err != nil {
		o.log.Error(ctx, err.Error())
		return err
	}

	err = o.redisCache.CancelOrder(ctx, reqParams)
	if err != nil {
		o.log.Error(ctx, err.Error())
	}
	return nil
}

func (o *OrderReceiver) FindAllOrders(ctx context.Context, reqParams *model.FindAllRequest) (*model.FindAllResponse, error) {
	//allOrders, err := o.redisCache.FindAllOrders(ctx, reqParams)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, order := range allOrders {
	//	fmt.Println(order)
	//}
	orders, paginationResponse, err := o.OrderRepository.FindAllOrders(ctx, reqParams)
	if err != nil {
		return nil, err
	}

	response := &model.FindAllResponse{
		Total:       paginationResponse.Total,
		CurrentPage: paginationResponse.CurrentPage,
		PerPage:     paginationResponse.PerPage,
		TotalInPage: paginationResponse.TotalInPage,
		LastPage:    paginationResponse.LastPage,
	}

	for _, order := range orders {
		orderResponse := &model.OrderResponse{
			OrderConsignmentID: order.OrderConsignmentID,
			OrderCreatedAt:     order.CreatedAt,
			OrderDescription:   order.ItemDescription,
			MerchantOrderID:    order.MerchantOrderID,
			RecipientName:      order.RecipientName,
			RecipientAddress:   order.RecipientAddress,
			RecipientPhone:     order.RecipientPhone,
			OrderAmount:        order.AmountToCollect,
			TotalFee:           order.TotalFee,
			Instruction:        order.SpecialInstruction,
			OrderTypeID:        order.OrderTypeID,
			CODFee:             order.CodFee,
			PromoDiscount:      order.PromoDiscount,
			Discount:           order.Discount,
			DeliveryFee:        order.DeliveryFee,
			OrderStatus:        order.OrderStatus.String(),
			OrderType:          order.DeliveryType.String(),
			ItemType:           order.ItemType.String(),
		}

		response.Orders = append(response.Orders, orderResponse)
	}

	return response, nil
}

// GenerateConsignmentID generates a unique consignment ID
func GenerateConsignmentID(prefix string, randomLength int) string {
	// Get the current date in YYMMDD format
	date := time.Now().Format("060102")

	// Generate a random alphanumeric suffix
	suffix := generateRandomAlphanumeric(randomLength)

	// Combine the prefix, date, and suffix
	return fmt.Sprintf("%s%s%s", prefix, date, suffix)
}

func generateRandomAlphanumeric(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		// Generate a secure random number
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic("failed to generate random alphanumeric string")
		}
		result[i] = charset[randomIndex.Int64()]
	}
	return string(result)
}
