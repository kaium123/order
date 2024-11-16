package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"

	"github.com/kaium123/order/internal/log"
	"github.com/labstack/echo/v4"
)

// JWTConfig holds the JWT key and other configurations
type JWTConfig struct {
	SecretKey string
}

// NewJWTMiddleware creates a new JWT middleware
func NewJWTMiddleware(config JWTConfig, log *log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the token from the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Error(c.Request().Context(), "Missing Authorization header")
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
			}

			// Check that the header starts with 'Bearer '
			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Error(c.Request().Context(), "Invalid token format")
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
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
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			// Validate the token
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// You can access the claims here if needed, e.g., claims["user_id"]
				// Set the claims in the context for later use
				c.Set("user_claims", claims)
				return next(c)
			} else {
				log.Error(c.Request().Context(), "Invalid token claims")
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
			}
		}
	}
}
