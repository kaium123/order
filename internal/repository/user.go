package repository

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// IUser defines the repository interface for user-related operations.
type IUser interface {
	FindUserByUserNameOrEmail(ctx context.Context, req *model.UserLoginRequest) (*model.User, error)
	SaveAccessToken(ctx context.Context, accessToken *model.AccessToken) error
	SaveRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) error
	RemoveAccessToken(ctx context.Context, userID int64) error
	RemoveRefreshToken(ctx context.Context, userID int64) error
}

type InitUserRepository struct {
	Db  *db.DB
	Log *log.Logger
}

type UserReceiver struct {
	log *log.Logger
	db  *db.DB
}

// NewUser creates a new instance of the User repository.
func NewUser(initUserRepository *InitUserRepository) IUser {
	return &UserReceiver{
		log: initUserRepository.Log,
		db:  initUserRepository.Db,
	}
}

// FindUserByUsername retrieves a user by username or email from the database.
func (u *UserReceiver) FindUserByUserNameOrEmail(ctx context.Context, req *model.UserLoginRequest) (*model.User, error) {
	user := &model.User{}
	err := u.db.NewSelect().
		Model(user).
		Where("email = ? OR user_name = ?", req.Email, req.Username).
		Limit(1).
		Scan(ctx)

	if err != nil {
		u.log.Error(ctx, fmt.Sprintf("Error finding user by username %s or email %s: %v", req.Username, req.Email, err))
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// SaveAccessToken saves a generated access token in the database.
func (u *UserReceiver) SaveAccessToken(ctx context.Context, accessToken *model.AccessToken) error {
	_, err := u.db.NewInsert().Model(accessToken).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save access token: %v", err)
	}
	return nil
}

// SaveRefreshToken saves a generated refresh token in the database.
func (u *UserReceiver) SaveRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) error {
	_, err := u.db.NewInsert().Model(refreshToken).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %v", err)
	}
	return nil
}

// RemoveAccessToken removes the access token from the database.
func (u *UserReceiver) RemoveAccessToken(ctx context.Context, userID int64) error {
	_, err := u.db.NewUpdate().Model(&model.AccessToken{}).
		Set("deleted_at = ?", time.Now().UTC()).
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		u.log.Error(ctx, fmt.Sprintf("Failed to remove access token for user %s: %v", userID, err))
		return fmt.Errorf("failed to remove access token")
	}
	return nil
}

// RemoveRefreshToken removes the refresh token from the database.
func (u *UserReceiver) RemoveRefreshToken(ctx context.Context, userID int64) error {
	_, err := u.db.NewUpdate().Model(&model.RefreshToken{}).
		Set("deleted_at = ?", time.Now().UTC()).
		Where("user_id = ?", userID).Exec(ctx)
	if err != nil {
		u.log.Error(ctx, fmt.Sprintf("Failed to remove refresh token for user %s: %v", userID, err))
		return fmt.Errorf("failed to remove refresh token")
	}
	return nil
}

// GenerateAccessToken creates and returns a new JWT access token for the user.
func (u *UserReceiver) GenerateAccessToken(userID string) (*model.AccessToken, error) {
	secretKey := "yourSecretKey"
	expiry := time.Now().Add(15 * time.Minute)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiry,
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	return &model.AccessToken{
		Token:     tokenString,
		UserID:    userID,
		Expiry:    expiry,
		CreatedAt: time.Now(),
	}, nil
}

// GenerateRefreshToken creates and returns a new refresh token for the user.
func (u *UserReceiver) GenerateRefreshToken(userID string) (*model.RefreshToken, error) {
	expiry := time.Now().Add(30 * 24 * time.Hour)

	refreshTokenString := uuid.New().String()

	return &model.RefreshToken{
		Token:     refreshTokenString,
		UserID:    userID,
		Expiry:    expiry,
		CreatedAt: time.Now(),
	}, nil
}

// ComparePassword compares the provided password with the stored hash value.
func ComparePassword(storedHash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match")
	}
	return nil
}
