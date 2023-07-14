package models

import "time"

type Wiki struct {
	ID          int
	Topic       string
	Description interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
