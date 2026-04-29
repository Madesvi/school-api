package models

import (
	"time"
)

// * указатели для null полей

type Exec struct {
	ID                  int        `gorm:"primaryKey" json:"id,omitempty"`
	FirstName           string     `json:"first_name,omitempty"`
	LastName            string     `json:"last_name,omitempty"`
	Email               string     `gorm:"uniqueIndex;not null" json:"email,omitempty"`
	UserName            string     `gorm:"column:username;uniqueIndex;not null" json:"username,omitempty"`
	Password            string     `json:"password,omitempty"`
	PasswordChangedAt   *time.Time `json:"password_changed_at,omitempty"`
	PasswordResetToken  *string    `gorm:"column:password_reset_token" json:"password_reset_token,omitempty"`
	PasswordCodeExpires *time.Time `json:"password_code_expires,omitempty"`
	Role                string     `json:"role,omitempty"`
	InactiveStatus      bool       `gorm:"column:inactive_status" json:"inactive_status,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UpdatePasswordResponse struct {
	Token          string `json:"token"`
	PasswordUpdate bool   `json:"password_update"`
}
