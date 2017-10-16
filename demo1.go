package main

import (
	"github.com/jinzhu/gorm"
	"github.com/geobe/go4oosy17/person"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	// Connect to database
	db := dbopen()
	// Close at program end
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&person.Person{})
	populate(db)

}

func dbopen() *gorm.DB {
	// Connect to database
	db, err := gorm.Open("postgres", "user=oosydev dbname=oosy17 password=myoosy17 sslmode=disable")
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return db
}

func populate(db *gorm.DB) {
	//var p Person
	var count int

	db.Model(&person.Person{}).Count(&count)
	if count == 0 {
		// Create
		db.Create(&person.Person{
			Lastname:  "Merkel",
			Firstname: "Angela",
			Username:  "angie",
			Password:  "Seehofer",
		})
		db.Create(&person.Person{
			Lastname:  "Seehofer",
			Firstname: "Horst",
			Username:  "hose",
			Password:  "lekreM",
		})
	}
}

