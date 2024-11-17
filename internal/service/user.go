// Package service provides the business logic for the User endpoint.
package service

import (
	"context"
	"errors"
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
		u.log.Error(ctx, err.Error())
		return nil, err
	}

	// Compare the provided password with the stored hash
	if !CheckPasswordHash(reqLogin.Password, user.PasswordHash) {
		u.log.Error(ctx, "password not matched")
		return nil, errors.New("password not matched")
	}

	err = u.Logout(ctx, user.ID)
	if err != nil {
		u.log.Error(ctx, err.Error())
		return nil, err
	}

	// Generate JWT tokens (access and refresh tokens)
	accessToken, err := u.jwtService.GenerateAccessToken(user.ID)
	if err != nil {
		u.log.Error(ctx, err.Error())
		return nil, err
	}

	refreshToken, err := u.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		u.log.Error(ctx, err.Error())
		return nil, err
	}

	// Save tokens to the database
	expirey := time.Now().Add(time.Hour * 1)
	err = u.UserRepository.SaveAccessToken(ctx, &model.AccessToken{
		Token:     accessToken,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		Expiry:    expirey, // Access token expires in 1 hour
	})
	if err != nil {
		u.log.Error(ctx, err.Error())
		return nil, err
	}

	err = u.UserRepository.SaveRefreshToken(ctx, &model.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		Expiry:    time.Now().Add(time.Hour * 24 * 7), // Refresh token expires in 7 days
	})
	if err != nil {
		u.log.Error(ctx, err.Error())
		return nil, err
	}

	// Store tokens in Redis
	key := fmt.Sprintf("access_token:%s", accessToken)
	err = u.redisCache.StoreToken(ctx, key, accessToken, 10*time.Minute)
	if err != nil {
		u.log.Error(ctx, err.Error())
	}

	key = fmt.Sprintf("refresh_token:%s", refreshToken)
	err = u.redisCache.StoreToken(ctx, key, refreshToken, 10*time.Minute)
	if err != nil {
		u.log.Error(ctx, err.Error())
	}

	// Return the utils with tokens
	return &model.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpireIn:     expirey.Unix(),
	}, nil
}

// Logout handles user logout and invalidates the tokens.
func (u *UserReceiver) Logout(ctx context.Context, userID int64) error {
	// Remove the access and refresh tokens from the database and cache
	accessTokens, err := u.UserRepository.RemoveAccessToken(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to remove access token: %v", err)
	}

	refreshTokens, err := u.UserRepository.RemoveRefreshToken(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to remove refresh token: %v", err)
	}

	// Optionally, invalidate session in Redis
	err = u.redisCache.InvalidateSession(ctx, userID)
	if err != nil {
		u.log.Error(ctx, fmt.Sprintf("Failed to invalidate session for user %s: %v", userID, err))
	}

	for _, accessToken := range accessTokens {
		key := fmt.Sprintf("access_token:%s", accessToken.Token)
		err := u.redisCache.DeleteKey(ctx, key)
		if err != nil {
			u.log.Error(ctx, fmt.Sprintf("Failed to invalidate session for user %s: %v", userID, err))
		}
	}
	for _, refreshToken := range refreshTokens {
		key := fmt.Sprintf("refresh_token:%s", refreshToken.Token)
		err := u.redisCache.DeleteKey(ctx, key)
		if err != nil {
			u.log.Error(ctx, fmt.Sprintf("Failed to invalidate session for user %s: %v", userID, err))
		}
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
