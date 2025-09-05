package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/configs"
)

func main() {
	// manual load ENV
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env \nCause: ", err.Error())
		return
	}

	// Inisialization databae for this project
	db, err := configs.InitDB()
	if err != nil {
		log.Println("Failed to connect to  database\nCause: ", err.Error())
		return
	}
	defer db.Close()

	// testing connection with database
	if err := configs.TestDB(db); err != nil {
		log.Println("Ping to DB failsed\nCause: ", err.Error())
	}
	log.Println("DB connected")

}
