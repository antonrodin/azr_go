package models

import (
	"time"
)

type Page struct {
	ID        int
	Published int
	Title     string
	Slug      string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
