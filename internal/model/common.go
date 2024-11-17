package model

import "fmt"

// ErrNotFound is the error for not found.
var ErrNotFound = fmt.Errorf("not found")

type PaginationResponse struct {
	Total       int `json:"total"`
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalInPage int `json:"total_in_page"`
	LastPage    int `json:"last_page"`
}
