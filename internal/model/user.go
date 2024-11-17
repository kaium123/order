package model

import (
	"errors"
	"github.com/uptrace/bun"
	"time"
)

// User represents a user in the system.
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID           int64     `json:"id" bun:"id,pk,autoincrement"`
	UserName     string    `json:"user_name" bun:"user_name"`
	Email        string    `json:"email" bun:"email"`
	PasswordHash string    `json:"-" bun:"password_hash"`
	CreatedAt    time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bun:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`
}

// AccessToken represents an access token for a user.
type AccessToken struct {
	bun.BaseModel `bun:"table:access_tokens"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Token     string    `json:"token" bun:"token"`
	UserID    int64     `json:"user_id" bun:"user_id"`
	Expiry    time.Time `json:"expiry" bun:"expiry"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	DeletedAt time.Time `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`
}

// RefreshToken represents a refresh token for re-authentication.
type RefreshToken struct {
	bun.BaseModel `bun:"table:refresh_tokens"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Token     string    `json:"token" bun:"token"`
	UserID    int64     `json:"user_id" bun:"user_id"`
	Expiry    time.Time `json:"expiry" bun:"expiry"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	DeletedAt time.Time `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`
}

// UserLoginRequest represents the login request payload.
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// UserLoginResponse represents the utils with tokens after a successful login.
type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpireIn     int64  `json:"expire_in"`
}

var ErrInvalidCredentials = errors.New("invalid credentials")
