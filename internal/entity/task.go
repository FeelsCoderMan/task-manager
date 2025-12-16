package entity

import "time"

type Task struct {
	ID          int       `json:"id",omitempty`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Priority    uint8     `json:"priority"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed,omitempty"`
}
