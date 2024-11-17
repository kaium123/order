package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/kaium123/order/internal/model"
	"github.com/kaium123/order/internal/repository"
	"github.com/kaium123/order/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/log"
	"github.com/labstack/echo/v4"
)

// JWTConfig holds the JWT key and other configurations
type JWTConfig struct {
	SecretKey  string
	DB         *db.DB // Add the database connection here
	RedisCache repository.IRedisCache
}

// NewJWTMiddleware creates a new JWT middleware
func NewJWTMiddleware(config JWTConfig, log *log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the token from the Authorization header
			var responseErr utils.ResponseError
			ctx := c.Request().Context()
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Error(ctx, "Missing Authorization header")
				return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, nil, "Unauthorized"))
			}

			// Check that the header starts with 'Bearer '
			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Error(ctx, "Invalid token format")
				return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, nil, "Unauthorized"))
			}

			// Extract the token from the Authorization header
			tokenString := authHeader[len("Bearer "):]

			// Parse and validate the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Ensure the token is signed with the correct method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
				}
				return []byte(config.SecretKey), nil
			})

			if err != nil {
				log.Error(ctx, fmt.Sprintf("Invalid token: %v", err))
				return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, nil, "Unauthorized"))
			}

			// Validate the token and extract the claims
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Extract the user_id from the claims
				if userID, ok := claims["user_id"].(float64); ok {
					key := fmt.Sprintf("access_token:%s", tokenString)
					fmt.Println("key : ", key)
					getToken, err := config.RedisCache.GetToken(ctx, key)
					if err != nil || getToken == "" {
						log.Error(ctx, "failed to found token from redis")
						// Check if the user exists and the token is valid in the database
						exists, err := config.DB.NewSelect().Model(&model.AccessToken{}).
							Where("user_id = ? and token = ? ", int64(userID), tokenString).
							Exists(ctx)

						if err != nil || !exists {
							log.Error(ctx, "invalid token or expired")
							return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, nil, "Unauthorized"))
						}
						err = config.RedisCache.StoreToken(ctx, key, tokenString, 10*time.Minute)
						if err != nil {
							log.Error(ctx, err.Error())
						}
					}

					fmt.Println("get token from redis - ", tokenString)

					// Set the user_id in the context
					c.Set("user_id", int64(userID))
				} else {
					log.Error(ctx, "User ID not found in token claims")
					return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, nil, "Unauthorized"))
				}

				// Optionally set all claims in the context if needed
				c.Set("user_claims", claims)

				// Proceed to the next handler
				return next(c)
			} else {
				log.Error(ctx, "Invalid token claims")
				return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, nil, "Unauthorized"))
			}
		}
	}
}
