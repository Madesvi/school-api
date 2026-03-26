package models

import "time"

type Exec struct {
	ID        int
	FirstName string `gorm:"not null" json:"first_name,omitempty"`
	LastName  string `gorm:"not null" json:"last_name,omitempty"`
	Email     string `gorm:"uniqueIndex;not null" json:"email,omitempty"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
