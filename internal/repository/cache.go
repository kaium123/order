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
	StoreToken(ctx context.Context, key string, token string, expiry time.Duration) error
	GetToken(ctx context.Context, key string) (string, error)
	DeleteKey(ctx context.Context, key string) error
	FindAllOrders(ctx context.Context, req *model.FindAllRequest) ([]model.Order, error)
}

type InitRedisCache struct {
	Client *redis.Client
	Log    *log.Logger
}

type redisCache struct {
	client *redis.Client
	log    *log.Logger
}

// NewRedisCache creates a new Redis client instance.
func NewRedisCache(initRedisCache *InitRedisCache) IRedisCache {
	return &redisCache{
		client: initRedisCache.Client,
		log:    initRedisCache.Log,
	}
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

// StoreAccessToken stores an access token in Redis with a specified expiration.
func (r *redisCache) StoreToken(ctx context.Context, key string, token string, expiry time.Duration) error {
	err := r.client.Set(ctx, key, token, expiry).Err()
	if err != nil {
		r.log.Error(ctx, fmt.Sprintf("Failed to store access token.", err))
		return fmt.Errorf("failed to store access token: %w", err)
	}
	return nil
}

// GetAccessToken retrieves an access token from Redis for a user.
func (r *redisCache) GetToken(ctx context.Context, key string) (string, error) {
	token, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		r.log.Error(ctx, fmt.Sprintf("Access token not found"))
		return "", nil
	} else if err != nil {
		r.log.Error(ctx, fmt.Sprintf("Failed to retrieve access token ", err))
		return "", fmt.Errorf("failed to retrieve access token: %w", err)
	}
	return token, nil
}

// DeleteKey removes a specific key from Redis.
func (r *redisCache) DeleteKey(ctx context.Context, key string) error {
	// Delete the specified key from Redis
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.log.Error(ctx, fmt.Sprintf("Failed to delete key %s from Redis: %v", key, err))
		return fmt.Errorf("failed to delete key: %w", err)
	}

	// Log the successful deletion for tracking
	r.log.Info(ctx, fmt.Sprintf("Successfully deleted key %s from Redis", key))
	return nil
}

// FindAllOrders retrieves orders from Redis based on the given filter and paginates them.
func (t *redisCache) FindAllOrders(ctx context.Context, req *model.FindAllRequest) ([]model.Order, error) {
	// Fetch order consignment IDs from the sorted set with score-based pagination (limit and offset)
	zRangeResp, err := t.client.ZRangeByScoreWithScores(ctx, "orders", &redis.ZRangeBy{
		Min:    "-inf", // Minimum score (start from the earliest order)
		Max:    "+inf", // Maximum score (end at the latest order)
		Offset: int64(req.Offset),
		Count:  int64(req.Limit),
	}).Result()

	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("Error fetching orders from sorted set: %v", err))
		return nil, err
	}

	// Initialize the list to store retrieved orders
	var orders []model.Order

	// For each order consignment ID, retrieve order details from the hash
	for _, z := range zRangeResp {
		orderKey := fmt.Sprintf("order:%s", z.Member)
		orderData, err := t.client.HGetAll(ctx, orderKey).Result()
		if err != nil {
			t.log.Error(ctx, fmt.Sprintf("Error fetching order data for consignment ID %s: %v", z.Member, err))
			continue
		}

		// Convert the hash data to an order object
		order := model.Order{
			OrderConsignmentID: orderData["consignment_id"],
			MerchantOrderID:    orderData["merchant_order_id"],
			//OrderStatus:        model.OrderStatus(orderData["order_status"]),
			CreatedAt: parseTime(orderData["created_at"]), // Use a helper function to parse time
			UpdatedAt: parseTime(orderData["updated_at"]),
		}

		// Apply filtering if an order status is specified
		if req.TransferStatus != "" && order.OrderStatus.String() != req.TransferStatus {
			continue // Skip this order if it doesn't match the filter
		}

		// Append the order to the result list
		orders = append(orders, order)
	}

	return orders, nil
}

// Helper function to parse time from a string
func parseTime(timeStr string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{} // Return zero value if parsing fails
	}
	return parsedTime
}
