package models

import (
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// type Person struct {
// 	ID        uint   `json:"id"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Email     string `json:"email"`
// 	Class     string `json:"class"`
// 	Subject   string `json:"subject"`
// }

type Person struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	Class     string `gorm:"not null" json:"class"`
	Subject   string `gorm:"not null" json:"subject"`
}

func (Person) TableName() string {
	return "teachers"
}

func GetPersonByIDFromPostgre(db *gorm.DB, id int) (*Person, error) {
	var person Person
	result := db.First(&person, id)
	if result.Error != nil {
		log.Error().Msg("Error")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	// fmt.Println("Person", person)
	return &person, nil
}
