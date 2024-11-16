package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"time"
)

// IRedisCache RedisListCache defines the interface for Redis list operations.
type IRedisCache interface {
	CacheOrder(ctx context.Context, order model.Order) error
	CancelOrder(ctx context.Context, reqParams *model.OrderCancelRequest) error
	InvalidateSession(ctx context.Context, userID int64) error
}

type InitRedisCache struct {
	Client *redis.Client
	Log    *log.Logger
}

type redisCache struct {
	client *redis.Client
	log    *log.Logger
}

// CacheOrder stores the order in Redis as a hash and adds it to a sorted set for easy retrieval.
func (t *redisCache) CacheOrder(ctx context.Context, order model.Order) error {
	orderKey := fmt.Sprintf("order:%s", order.OrderConsignmentID) // Use order ID as key

	// Convert time to string (ISO 8601 or Unix timestamp)
	createdAtStr := order.CreatedAt.Format(time.RFC3339) // Using ISO 8601 format
	updatedAtStr := order.UpdatedAt.Format(time.RFC3339)

	// Store the order details as a Redis hash
	_, err := t.client.HSet(ctx, orderKey, map[string]interface{}{
		"consignment_id":    order.OrderConsignmentID,
		"merchant_order_id": order.MerchantOrderID,
		"order_status":      order.OrderStatus,
		"delivery_fee":      order.DeliveryFee,
		"created_at":        createdAtStr,
		"updated_at":        updatedAtStr,
	}).Result()

	if err != nil {
		// Log the error for debugging
		t.log.Error(ctx, fmt.Sprintf("Error caching order with ID %s: %v", order.OrderConsignmentID, err))
		return err
	}

	// Add the order to the sorted set with CreatedAt as the score
	zAddErr := t.client.ZAdd(ctx, "orders", &redis.Z{
		Score:  float64(order.CreatedAt.Unix()),
		Member: order.OrderConsignmentID,
	}).Err()

	if zAddErr != nil {
		// Log the error for debugging
		t.log.Error(ctx, fmt.Sprintf("Error adding order to sorted set for order ID %s: %v", order.OrderConsignmentID, zAddErr))
		return zAddErr
	}

	return nil
}

// CancelOrder invalidates the order cache and removes the order from the sorted set in Redis.
func (t *redisCache) CancelOrder(ctx context.Context, reqParams *model.OrderCancelRequest) error {
	// Invalidate the order cache in Redis (delete the order hash)
	orderKey := fmt.Sprintf("order:%s", reqParams.ConsignmentID)
	_, err := t.client.Del(ctx, orderKey).Result()
	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("Failed to invalidate cache for order with ID %s: %v", reqParams.ConsignmentID, err))
		return err
	}

	// Remove the order from the sorted set (remove the member from the 'orders' sorted set)
	_, err = t.client.ZRem(ctx, "orders", reqParams.ConsignmentID).Result()
	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("Failed to remove order from sorted set for order with ID %s: %v", reqParams.ConsignmentID, err))
		return err
	}

	return nil
}

// InvalidateSession removes the session (such as JWT token) for the specified user.
func (r *redisCache) InvalidateSession(ctx context.Context, userID int64) error {
	// Redis key for storing user session
	sessionKey := fmt.Sprintf("session:%d", userID)

	// Remove the session key from Redis
	err := r.client.Del(ctx, sessionKey).Err()
	if err != nil {
		// Log the error if invalidating the session fails
		r.log.Error(ctx, fmt.Sprintf("Failed to invalidate session for user %s: %v", userID, err))
		return fmt.Errorf("failed to invalidate session")
	}

	// If the session was invalidated successfully, return nil
	return nil
}

// NewRedisCache creates a new Redis client instance.
func NewRedisCache(initRedisCache *InitRedisCache) IRedisCache {
	return &redisCache{
		client: initRedisCache.Client,
		log:    initRedisCache.Log,
	}
}
