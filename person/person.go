package person

import "github.com/jinzhu/gorm"

type Person struct {
	gorm.Model
//	ID        uint `gorm:"primary_key"`
	Lastname  string
	Firstname string
	Username  string
	Password  string
}

