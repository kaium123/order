package service

import (
	"github.com/golang-jwt/jwt"
	"time"
)

// IJWTService defines the methods that our JWT service should implement.
type IJWTService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
}

// JWTService is the concrete implementation of the IJWTService interface.
type JWTService struct {
	secretKey string
}

// NewJWTService creates a new JWTService instance with the given secret key.
func NewJWTService(secretKey string) IJWTService {
	return &JWTService{secretKey: secretKey}
}

// GenerateAccessToken generates an access token for the user.
func (s *JWTService) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,                               // Subject is the user ID
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // Expiry time set to 1 hour
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// GenerateRefreshToken generates a refresh token for the user.
func (s *JWTService) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,                                    // Subject is the user ID
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // Expiry time set to 7 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}
