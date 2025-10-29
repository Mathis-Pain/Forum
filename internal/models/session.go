package models

import "time"

type Session struct {
	ID        string                 `json:"id"`
	UserID    int                    `json:"user_id"`
	Data      map[string]interface{} `json:"data"`
	ExpiresAt time.Time              `json:"expires_at"`
	CreatedAt time.Time              `json:"created_at"`
}
