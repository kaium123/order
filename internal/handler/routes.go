package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/middleware"
	"github.com/kaium123/order/internal/repository"
	"github.com/kaium123/order/internal/service"
	"github.com/labstack/echo/v4"
)

type ServiceRegistry struct {
	EchoEngine  *echo.Echo
	RedisClient *redis.Client
	DBInstance  *db.DB
	Log         *log.Logger
}

// Register registers the routes for the application.
func Register(serviceRegistry *ServiceRegistry) {
	serviceRegistry.EchoEngine.Validator = &CustomValidator{validator: validator.New()}

	api := serviceRegistry.EchoEngine.Group("/api/v1")

	// Health check
	healthHandler := NewHealth()
	api.GET("/healthz", healthHandler.Healthz)

	// Inject Order Dependency
	redisRepository := repository.NewRedisCache(&repository.InitRedisCache{
		Client: serviceRegistry.RedisClient,
		Log:    serviceRegistry.Log,
	})

	// Initialize JWT middleware
	jwtMiddleware := middleware.NewJWTMiddleware(middleware.JWTConfig{
		SecretKey:  "123",
		DB:         serviceRegistry.DBInstance,
		RedisCache: redisRepository,
	}, serviceRegistry.Log)

	orderRepository := repository.NewOrder(&repository.InitOrderRepository{
		Db: serviceRegistry.DBInstance, Log: serviceRegistry.Log,
	})
	orderService := service.NewOrder(&service.InitOrderService{
		Log: serviceRegistry.Log, OrderRepository: orderRepository, RedisCache: redisRepository,
	})
	orderHandler := NewOrder(&InitOrderHandler{
		Service: orderService, Log: serviceRegistry.Log,
	})

	// Inject Auth Dependency
	userRepository := repository.NewUser(&repository.InitUserRepository{
		Db: serviceRegistry.DBInstance, Log: serviceRegistry.Log,
	})
	jwtService := service.NewJWTService("123")
	authService := service.NewUser(&service.InitUserService{
		Log: serviceRegistry.Log, UserRepository: userRepository,
		RedisCache: redisRepository,
		JWTService: jwtService,
	})
	authHandler := NewAuth(&InitAuthHandler{
		Service: authService, Log: serviceRegistry.Log,
	})

	// Add routes for order
	order := api.Group("/orders", jwtMiddleware)
	{
		order.POST("", orderHandler.CreateOrder)
		order.GET("/all", orderHandler.FindAllOrders)
		order.PUT("/:CONSIGNMENT_ID/cancel", orderHandler.CancelOrder)
	}

	// Add routes for auth (login and logout)
	api.POST("/login", authHandler.Login)
	api.POST("/logout", authHandler.Logout, jwtMiddleware)

}
