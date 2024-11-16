package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/kaium123/order/internal/model"
	"github.com/kaium123/order/internal/utils"
	"net/http"
	"strings"

	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/log"
	"github.com/labstack/echo/v4"
)

// JWTConfig holds the JWT key and other configurations
type JWTConfig struct {
	SecretKey string
	DB        *db.DB // Add the database connection here
}

// NewJWTMiddleware creates a new JWT middleware
func NewJWTMiddleware(config JWTConfig, log *log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the token from the Authorization header
			var responseErr utils.ResponseError
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Error(c.Request().Context(), "Missing Authorization header")
				return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, nil, "Unauthorized"))
			}

			// Check that the header starts with 'Bearer '
			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Error(c.Request().Context(), "Invalid token format")
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
				log.Error(c.Request().Context(), fmt.Sprintf("Invalid token: %v", err))
				return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, map[string][]string{"token_parsing_error": []string{err.Error()}}, "Unauthorized"))
			}

			// Validate the token and extract the claims
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Extract the user_id from the claims
				if userID, ok := claims["user_id"].(string); ok {
					// Check if the user exists and the token is valid in the database
					exists, err := config.DB.NewSelect().Model(&model.AccessToken{}).
						Where("user_id = ? and token = ? ", userID, tokenString).
						Exists(context.Background())

					if err != nil || !exists {
						log.Error(context.Background(), "invalid token or expired")
						return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, map[string][]string{"token_validation_error": []string{err.Error()}}, "Unauthorized"))
					}

					// Set the user_id in the context
					c.Set("user_id", userID)
				} else {
					log.Error(c.Request().Context(), "User ID not found in token claims")
					return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, map[string][]string{"token_validation_error": []string{err.Error()}}, "Unauthorized"))
				}

				// Optionally set all claims in the context if needed
				c.Set("user_claims", claims)

				// Proceed to the next handler
				return next(c)
			} else {
				log.Error(c.Request().Context(), "Invalid token claims")
				return c.JSON(responseErr.GetErrorResponse(http.StatusUnauthorized, map[string][]string{"token_validation_error": []string{err.Error()}}, "Unauthorized"))
			}
		}
	}
}
