package model

import (
	"errors"
	"time"
)

// User represents a user in the system.
type User struct {
	ID           string    `json:"id" bun:"id"`                 // Unique identifier for the user
	UserName     string    `json:"user_name" bun:"user_name"`   // Username chosen by the user
	Email        string    `json:"email" bun:"email"`           // User's email address
	PasswordHash string    `json:"-" bun:"password_hash"`       // Hashed password (never store plain text password)
	CreatedAt    time.Time `json:"created_at" bun:"created_at"` // Timestamp for when the user was created
	UpdatedAt    time.Time `json:"updated_at" bun:"updated_at"` // Timestamp for when the user was last updated
}

// AccessToken represents an access token for a user.
type AccessToken struct {
	Token     string    `json:"token" bun:"token"`           // JWT or token string
	UserID    string    `json:"user_id" bun:"user_id"`       // The ID of the user to whom the token belongs
	Expiry    time.Time `json:"expiry" bun:"expiry"`         // The expiration timestamp of the access token
	CreatedAt time.Time `json:"created_at" bun:"created_at"` // The timestamp when the token was created
}

// RefreshToken represents a refresh token for re-authentication.
type RefreshToken struct {
	Token     string    `json:"token" bun:"token"`           // Refresh token string
	UserID    string    `json:"user_id" bun:"user_id"`       // The ID of the user to whom the token belongs
	Expiry    time.Time `json:"expiry" bun:"expiry"`         // Expiration timestamp for the refresh token
	CreatedAt time.Time `json:"created_at" bun:"created_at"` // The timestamp when the refresh token was created
}

// UserLoginRequest represents the login request payload.
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// UserLoginResponse represents the response with tokens after a successful login.
type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

var ErrInvalidCredentials = errors.New("invalid credentials")
