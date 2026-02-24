package models

import (
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"int"`
}

func (Person) TableName() string {
	return "my_table_db"
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
