// Package models use for define models for db
package models

import "time"

// type Teacher struct {
// 	ID        int    `json:"id,omitempty"`         // 8 bytes
// 	FirstName string `json:"first_name,omitempty"` // 16 bytes
// 	LastName  string `json:"last_name,omitempty"`  // 16 bytes
// 	Class     string `json:"class,omitempty"`      // 16 bytes
// 	Subject   string `json:"subject,omitempty"`    // 16 bytesk
// 	// 74 bytes < 128 bytes for value
// }

type Teacher struct {
	ID        int    `gorm:"primaryKey" json:"id,omitempty"`
	FirstName string `gorm:"not null" json:"first_name,omitempty"`
	LastName  string `gorm:"not null" json:"last_name,omitempty"`
	Email     string `gorm:"uniqueIndex;not null" json:"email,omitempty"`
	Class     string `gorm:"not null" json:"class,omitempty"`
	Subject   string `gorm:"not null" json:"subject,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
