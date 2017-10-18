package main

import (
	"github.com/geobe/go4oosy17/person"
)

func main() {
	// Connect to database
	db := person.Dbopen()
	// Close at program end
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&person.Person{})
	person.Populate(db)

}

