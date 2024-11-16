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
	consignmentID := GenerateConsignmentID("CN", 6)
	reqOrder.OrderConsignmentID = consignmentID

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
		OrderStatus:     order.OrderStatus,
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
	orders, err := o.OrderRepository.FindAllOrders(ctx, reqParams)
	if err != nil {
		return nil, err
	}

	return orders, nil
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
