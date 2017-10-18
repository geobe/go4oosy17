package person

import "github.com/jinzhu/gorm"

type Person struct {
	gorm.Model
	Lastname  string
	Firstname string
	Username  string	`gorm:"not null;unique"`
	Password  string
}

