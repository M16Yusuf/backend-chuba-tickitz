package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// manual load ENV
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env \nCause: ", err.Error())
		return
	}
}
