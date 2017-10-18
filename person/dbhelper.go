package person

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Dbopen() *gorm.DB {
	// Connect to database
	db, err := gorm.Open("postgres", "user=oosydev dbname=oosy17 password=myoosy17 sslmode=disable")
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return db
}

func Populate(db *gorm.DB) {
	//var p Person
	var count int

	db.Model(&Person{}).Count(&count)
	if count == 0 {
		// Create
		db.Create(&Person{
			Lastname:  "Merkel",
			Firstname: "Angela",
			Username:  "angie",
			Password:  "Seehofer",
		})
		db.Create(&Person{
			Lastname:  "Seehofer",
			Firstname: "Horst",
			Username:  "hose",
			Password:  "lekreM",
		})
	}
}
