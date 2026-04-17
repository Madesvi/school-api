package models

import (
	"time"
)

// * указатели для null полей

type Exec struct {
	ID                  int        `gorm:"primaryKey"`
	FirstName           string     `json:"first_name,omitempty"`
	LastName            string     `json:"last_name,omitempty"`
	Email               string     `gorm:"uniqueIndex;not null" json:"email"`
	Username            string     `gorm:"uniqueIndex;not null" json:"username"`
	Password            string     `json:"-"`
	PasswordChangedAt   *time.Time `json:"password_changed_at,omitempty"`
	PasswordResetCode   *string    `gorm:"column:password_reset_token" json:"password_reset_token,omitempty"`
	PasswordCodeExpires *time.Time `json:"password_code_expires,omitempty"`
	Role                string     `json:"role"`
	InactiveStatus      bool       `gorm:"column:inactive_status" json:"inactive_status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
