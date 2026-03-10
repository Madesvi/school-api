// Package models use for define models for db
package models

import "time"

type Teacher struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	FirstName string `gorm:"not null" json:"first_name,omitempty"`
	LastName  string `gorm:"not null" json:"last_name,omitempty"`
	Email     string `gorm:"uniqueIndex;not null" json:"email,omitempty"`
	Class     string `gorm:"not null" json:"class,omitempty"`
	Subject   string `gorm:"not null" json:"subject,omitempty"`

	Students []Student `gorm:"foreignKey:TeacherID" json:"students,omitempty"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
