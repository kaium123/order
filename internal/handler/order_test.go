package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	cache2 "github.com/kaium123/order/internal/cache"
	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"github.com/kaium123/order/internal/repository"
	"github.com/kaium123/order/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func InitSetup(t *testing.T) TodoHandler {
	logger := log.New()
	dbInstance, err := db.NewMemory()
	require.NoError(t, err)
	err = db.Migrate(dbInstance)
	require.NoError(t, err)

	redis := cache2.New(&cache2.Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       5,
	})

	redisRepository := repository.NewRedisCache(&repository.InitRedisCache{
		Client: redis, Log: logger,
	})
	redisRepository.DeleteAll(context.Background())
	repository := repository.NewTodo(&repository.InitTodoRepository{Db: dbInstance, Log: logger})
	service := service.NewTodo(&service.InitTodoService{Log: logger, TodoRepository: repository, RedisCache: redisRepository})
	todoHandler := NewTodo(&InitTodoHandler{Service: service, Log: logger})
	return todoHandler
}

func TestTodoHandler_Create(t *testing.T) {
	type want struct {
		StatusCode int
		Response   []byte
	}
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	handler := InitSetup(t)

	tests := []struct {
		name       string
		createBody string
		want       want
		wantErr    bool
	}{
		{
			name:       "successful_create",
			createBody: `{"task":"Created Task","priority":"high"}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":{"Task":"Created Task","Status":"created","Priority":"high"}}`),
			},
		},
		{
			name:       "successful_create_but_with_ignore_status",
			createBody: `{"task":"Created Task", "status":"done","priority":"high"}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":{"Task":"Created Task","Status":"created","Priority":"high","ID":1}}`), // Excluded timestamps
			},
		},
		{
			name:       "invalid_request_body",
			createBody: `{"task":1}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare
			req := httptest.NewRequest(http.MethodPost, "/dummy/target", bytes.NewReader([]byte(tt.createBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/todos")

			// Execute
			require.NoError(t, handler.Create(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			// Compare ignoring CreatedAt and UpdatedAt fields
			opts := []cmp.Option{
				cmpTransformJSON(t),
				ignoreMapEntires(map[string]any{"CreatedAt": 1, "UpdatedAt": 1, "ID": 1}), // Ignore these fields
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func TestTodoHandler_Update(t *testing.T) {
	type want struct {
		StatusCode int
		Response   []byte
	}

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	handler := InitSetup(t)

	tests := []struct {
		name       string
		createBody string
		updateBody string
		updateID   string
		want       want
		wantErr    bool
	}{
		{
			name:       "successful_update",
			createBody: `{"task":"Updated Task","priority":"high"}`,
			updateBody: `{"task":"Updated Task","status":"done","priority":"high"}`,
			want: want{
				StatusCode: http.StatusOK,
				Response:   []byte(`{"data":{"Task":"Updated Task","Status":"done"}}`), // Only Task and Status
			},
		},
		{
			name:       "not_found_record",
			updateID:   "-1",
			updateBody: `{"task":"Updated Task","status":"done"}`,
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name:       "invalid_request_body",
			updateBody: `{"task":1}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "invalid_request_parameter",
			updateID:   "invalid",
			updateBody: `{"task":"Updated Task","status":"done"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare
			id := "1"
			if tt.updateID != "" {
				id = tt.updateID
			} else if tt.createBody != "" {
				id = strconv.Itoa(createTask(t, e, handler, tt.createBody))
			}

			req := httptest.NewRequest(http.MethodPut, "/dummy/target", bytes.NewReader([]byte(tt.updateBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/todos/:id")
			c.SetParamNames("id")
			c.SetParamValues(id)

			// Execute
			require.NoError(t, handler.Update(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			// Compare ignoring CreatedAt, UpdatedAt, ID, and Priority fields
			opts := []cmp.Option{
				cmpTransformJSON(t),
				ignoreMapEntires(map[string]any{"CreatedAt": 1, "UpdatedAt": 1, "ID": 1, "Priority": 1}), // Exclude these fields
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func TestTodoHandler_Delete(t *testing.T) {
	type want struct {
		StatusCode int
	}

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	handler := InitSetup(t)

	tests := []struct {
		name       string
		createBody string
		deleteID   string
		want       want
		wantErr    bool
	}{
		{
			name:       "successful_delete",
			createBody: `{"task":"Deleted Task","priority":"high"}`,
			want: want{
				StatusCode: http.StatusNoContent,
			},
		},
		{
			name:     "not_found_record",
			deleteID: "-1",
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name:     "invalid_request_parameter",
			deleteID: "invalid",
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare
			id := "1"
			if tt.deleteID != "" {
				id = tt.deleteID
			} else if tt.createBody != "" {
				id = strconv.Itoa(createTask(t, e, handler, tt.createBody))
			}

			req := httptest.NewRequest(http.MethodDelete, "/dummy/target", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/todos/:id")
			c.SetParamNames("id")
			c.SetParamValues(id)

			// Execute
			require.NoError(t, handler.Delete(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)
		})
	}
}

func TestTodoHandler_Find(t *testing.T) {
	type want struct {
		StatusCode int
		Response   []byte
	}

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	handler := InitSetup(t)

	tests := []struct {
		name       string
		createBody string
		findID     string
		want       want
		wantErr    bool
	}{
		{
			name:       "successful_find",
			createBody: `{"task":"Found Task","priority":"high"}`,
			want: want{
				StatusCode: http.StatusOK,
				Response:   []byte(`{"data":{"Task":"Found Task","Status":"created","Priority":"high"}}`),
			},
		},
		{
			name:   "not_found",
			findID: "-1",
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name:   "invalid_request_parameter",
			findID: "invalid",
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare
			id := "1"
			if tt.findID != "" {
				id = tt.findID
			} else if tt.createBody != "" {
				id = strconv.Itoa(createTask(t, e, handler, tt.createBody))
			}

			req := httptest.NewRequest(http.MethodGet, "/dummy/target", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/todos/:id")
			c.SetParamNames("id")
			c.SetParamValues(id)

			// Execute
			require.NoError(t, handler.Find(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			opts := []cmp.Option{
				cmpTransformJSON(t),
				ignoreMapEntires(map[string]any{"CreatedAt": 1, "UpdatedAt": 1, "ID": 1}),
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func TestTodoHandler_FindAll(t *testing.T) {
	type want struct {
		StatusCode int
		Response   []map[string]any // Changed to a slice of maps to handle dynamic fields
	}

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	handler := InitSetup(t)

	tests := []struct {
		name       string
		createBody []string
		want       want
		wantErr    bool
	}{
		{
			name:       "successful_find_all",
			createBody: []string{`{"task":"Task A", "priority":"low"}`, `{"task":"Task B", "priority":"medium"}`},
			want: want{
				StatusCode: http.StatusOK,
				Response: []map[string]any{
					{"Task": "Task A", "Priority": "low"},
					{"Task": "Task B", "Priority": "medium"},
				},
			},
		},
		{
			name: "successful_find_all_but_empty",
			want: want{
				StatusCode: http.StatusOK,
				Response:   []map[string]any{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create tasks for the test
			for _, body := range tt.createBody {
				createTask(t, e, handler, body)
			}

			req := httptest.NewRequest(http.MethodGet, "/dummy/target", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/todos")

			// Execute
			require.NoError(t, handler.FindAll(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}

			// Parse the response body into a map
			var gotResponse map[string]any
			err := json.Unmarshal(rec.Body.Bytes(), &gotResponse)
			require.NoError(t, err)
		})
	}
}

func createTask(t *testing.T, e *echo.Echo, handler TodoHandler, body string) int {
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader([]byte(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Create(c)
	require.NoError(t, err)

	type response struct {
		Data   model.Todo
		Status string
	}

	var res response
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err, "Failed to json.Unmarshal err: %s", err)
	require.NotEmpty(t, res.Data.ID)

	return res.Data.ID
}
