package models

type Teacher struct {
	ID        int    `json:"id,omitempty"`         // 8 bytes
	FirstName string `json:"first_name,omitempty"` // 16 bytes
	LastName  string `json:"last_name,omitempty"`  // 16 bytes
	Class     string `json:"class,omitempty"`      // 16 bytes
	Subject   string `json:"subject,omitempty"`    // 16 bytesk
	// 74 bytes < 128 bytes for value
}
