package repository

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"go.uber.org/zap"

	"strconv"
	"time"
)

// IRedisCache RedisListCache defines the interface for Redis list operations.
type IRedisCache interface {
	Get(ctx context.Context, key string) (string, error)
	Add(ctx context.Context, todoKey string, todo *model.Todo) error
	FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error)
	DeleteAll(ctx context.Context) error
	Delete(ctx context.Context, todoKey string) error
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

// Get retrieves a value from Redis using a key.
func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("key does not exist: %s", key)
	} else if err != nil {
		return "", err
	}
	return val, nil
}

// Delete removes a key from Redis if it exists
func (td *redisCache) Delete(ctx context.Context, todoKey string) error {
	// Execute the DEL command to remove the key from Redis
	err := td.client.Del(ctx, todoKey).Err()
	if err != nil {
		td.log.Error(ctx, "Failed to delete key from Redis", zap.Error(err))
		return err
	}

	td.log.Info(ctx, "Successfully deleted key from Redis", zap.String("key", todoKey))
	return nil
}

func (td *redisCache) Add(ctx context.Context, todoKey string, todo *model.Todo) error {
	// Encode ID as Base64 before storing
	encodedID := base64.StdEncoding.EncodeToString([]byte(todoKey))

	// Store the Todo details
	err := td.client.HMSet(ctx, todoKey, map[string]interface{}{
		"Id":        strconv.Itoa(todo.ID),
		"Task":      todo.Task,
		"Status":    string(todo.Status),
		"Priority":  string(todo.Priority),
		"CreatedAt": todo.CreatedAt,
		"UpdatedAt": todo.UpdatedAt,
	}).Err()
	if err != nil {
		td.log.Error(ctx, "in hset ", zap.Error(err))
		return err
	}

	// Calculate the score based on your sorting criteria
	score := CalculateScore(todo)

	fmt.Println("task ", todo.Task, "score --  ", score)

	// Add to sorted set for ordering
	err = td.client.ZAdd(ctx, "todos_sorted", &redis.Z{
		Score:  float64(score),
		Member: encodedID,
	}).Err()
	if err != nil {
		td.log.Error(ctx, "Error flushing database", zap.Error(err))
		return err
	}

	return nil
}

func (td *redisCache) DeleteAll(ctx context.Context) error {
	// Use FLUSHDB to remove all keys in the current database
	_, err := td.client.FlushDB(ctx).Result()
	if err != nil {
		td.log.Error(ctx, "Error flushing database", zap.Error(err))
		return err
	}

	return nil
}

func CalculateScore(todo *model.Todo) float64 {
	// Calculate score based on status, priority, and timestamps
	var score float64

	// Step 1: Lower score for 'done' status
	if todo.Status == "done" {
		score = 0 // Set score to the lowest for done tasks
	} else {
		// Step 2: Use created_at timestamp for sorting non-done tasks
		score += (float64(todo.CreatedAt.Unix()) / 100000)
		fmt.Println(score)

		// Step 3: Adjust score based on priority
		switch todo.Priority {
		case "high":
			score = score * 3.0 // High priority gets the highest score
		case "medium":
			score = score * 2.0 // Medium priority gets a medium score
		case "low":
			score = score*1.0 + 1.0 // Low priority gets the lowest score
		}
	}

	return score
}

func (td *redisCache) FindAll(ctx context.Context, reqParams *model.FindAllRequest) ([]*model.Todo, error) {
	var todos []*model.Todo

	// Retrieve IDs from the Sorted Set based on your sorting criteria
	todoIDs, err := td.client.ZRevRange(ctx, "todos_sorted", 0, -1).Result()
	if err != nil {
		td.log.Error(ctx, "zrev range ", zap.Error(err))
		return nil, err
	}

	if len(todoIDs) == 0 {
		return nil, errors.New("no todos found")
	}

	// Fetching Todos from Redis
	for _, id := range todoIDs {
		d, err := base64.StdEncoding.DecodeString(id)
		//fmt.Println(string(d))
		if err != nil {
			continue
		}

		// ID must be encoded in Base64; decode it
		todoKey := string(d)

		// Fetch all fields of the todo item
		todoData, err := td.client.HGetAll(ctx, todoKey).Result()
		if err != nil {
			td.log.Error(ctx, "Error fetching todo from Redis", zap.String("todoKey", todoKey), zap.Error(err))
			continue
		}

		todoIDInt, _ := strconv.Atoi(todoData["Id"])
		todo := &model.Todo{
			ID:        todoIDInt,
			Task:      todoData["Task"],
			Status:    model.Status(todoData["Status"]),
			Priority:  model.TodoPriority(todoData["Priority"]),
			CreatedAt: parseTime(todoData["CreatedAt"]),
			UpdatedAt: parseTime(todoData["UpdatedAt"]),
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func parseTime(timeStr string) time.Time {
	layout := time.RFC3339
	t, _ := time.Parse(layout, timeStr)
	return t
}
