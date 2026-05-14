package models

import "time"

type Student struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	FirstName string `gorm:"not null" json:"first_name,omitempty"`
	LastName  string `gorm:"not null" json:"last_name,omitempty"`
	Email     string `gorm:"uniqueIndex;not null" json:"email,omitempty"`
	TeacherID int    `gorm:"not null" json:"teacher_id,omitempty"`
	Balance   int64  `gorm:"not null" json:"balance,omitempty"`

	Teacher *Teacher `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
