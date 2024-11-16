package handler

import (
	cache2 "github.com/kaium123/order/internal/cache"
	"github.com/kaium123/order/internal/db"

	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	// Setup
	e := echo.New()
	cache := cache2.New(&cache2.Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       5,
	})
	dbInstance, err := db.NewMemory()
	require.NoError(t, err)
	err = db.Migrate(dbInstance)
	require.NoError(t, err)
	logger := log.New()
	Register(&ServiceRegistry{
		EchoEngine:  e,
		DBInstance:  dbInstance,
		Log:         logger,
		RedisClient: cache,
	})

	// Test cases
	tests := []struct {
		name         string
		method       string
		target       string
		expectedCode int
	}{
		{"Health_Check", http.MethodGet, "/api/v1/healthz", http.StatusOK},
		{"Create_Todo_without_body", http.MethodPost, "/api/v1/todos", http.StatusBadRequest}, // Assuming no body is sent, should return BadRequest
		{"Get_all_Todos", http.MethodGet, "/api/v1/todos", http.StatusOK},
		{"Get_non-existent_Todo", http.MethodGet, "/api/v1/todos/1", http.StatusNotFound},       // Assuming no todo with id 1 exists
		{"Update_Todo_without_body", http.MethodPut, "/api/v1/todos/1", http.StatusNotFound},    // Assuming no body is sent, should return BadRequest
		{"Delete_non-existent_Todo", http.MethodDelete, "/api/v1/todos/1", http.StatusNotFound}, // Assuming no todo with id 1 exists
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.target, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
