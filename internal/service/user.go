// Package service provides the business logic for the User endpoint.
package service

import (
	"context"
	"fmt"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"github.com/kaium123/order/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// IUser is the service interface for user-related operations.
type IAuth interface {
	Login(ctx context.Context, reqLogin *model.UserLoginRequest) (*model.UserLoginResponse, error)
	Logout(ctx context.Context, userID int64) error
}

type UserReceiver struct {
	log            *log.Logger
	UserRepository repository.IUser
	redisCache     repository.IRedisCache
	jwtService     IJWTService
}

type InitUserService struct {
	Log            *log.Logger
	UserRepository repository.IUser
	RedisCache     repository.IRedisCache
	JWTService     IJWTService
}

// NewUser creates a new User service.
func NewUser(initUserService *InitUserService) IAuth {
	return &UserReceiver{
		log:            initUserService.Log,
		UserRepository: initUserService.UserRepository,
		redisCache:     initUserService.RedisCache,
		jwtService:     initUserService.JWTService,
	}
}

// Login handles user login, generates JWT tokens, and saves them.
func (u *UserReceiver) Login(ctx context.Context, reqLogin *model.UserLoginRequest) (*model.UserLoginResponse, error) {
	// Validate user credentials
	user, err := u.UserRepository.FindUserByUserNameOrEmail(ctx, reqLogin)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Compare the provided password with the stored hash
	if !CheckPasswordHash(reqLogin.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT tokens (access and refresh tokens)
	accessToken, err := u.jwtService.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := u.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	// Save tokens to the database
	err = u.UserRepository.SaveAccessToken(ctx, &model.AccessToken{
		Token:     accessToken,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		Expiry:    time.Now().Add(time.Hour * 1), // Access token expires in 1 hour
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save access token: %v", err)
	}

	err = u.UserRepository.SaveRefreshToken(ctx, &model.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		Expiry:    time.Now().Add(time.Hour * 24 * 7), // Refresh token expires in 7 days
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %v", err)
	}

	// Return the response with tokens
	return &model.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
	}, nil
}

// Logout handles user logout and invalidates the tokens.
func (u *UserReceiver) Logout(ctx context.Context, userID int64) error {
	// Remove the access and refresh tokens from the database and cache
	err := u.UserRepository.RemoveAccessToken(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to remove access token: %v", err)
	}

	err = u.UserRepository.RemoveRefreshToken(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to remove refresh token: %v", err)
	}

	// Optionally, invalidate session in Redis
	err = u.redisCache.InvalidateSession(ctx, userID)
	if err != nil {
		u.log.Error(ctx, fmt.Sprintf("Failed to invalidate session for user %s: %v", userID, err))
	}

	return nil
}

// CheckPasswordHash compares the provided password with the stored hash.
func CheckPasswordHash(providedPassword, storedHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedPassword))
	if err != nil {
		// If the password doesn't match, return false
		return false
	}
	// If the password matches the hash, return true
	return true
}
