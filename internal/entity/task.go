package entity

import "time"

type Task struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Priority    uint8     `json:"priority,omitempty"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed,omitempty"`
}
