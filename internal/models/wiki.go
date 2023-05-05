package models

import "time"

type Wiki struct {
	ID          int
	Topic       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
