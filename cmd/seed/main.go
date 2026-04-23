package main

import (
	"log"

	"github.com/umanchanda/NBA-API/database"
)

func main() {
	if err := database.DatabaseFunctions(); err != nil {
		log.Fatal(err)
	}
	log.Println("database seeded successfully")
}
