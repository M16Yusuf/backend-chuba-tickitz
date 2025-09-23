package main

import (
	"context"
	"log"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/configs"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/routers"
)

// @title 											CHUBA TICKITZ
// @version 										1.0
// @description 								Ticketing application with restful API powered by gin
// @host												127.0.0.1:6969/api/
// @securityDefinitions.apikey 	JWTtoken
// @in header
// @name Authorization
func main() {
	// manual load ENV
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env \nCause: ", err.Error())
		// return
	}

	// Inisialization databae for this project
	db, err := configs.InitDB()
	if err != nil {
		log.Println("Failed to connect to  database\nCause: ", err.Error())
		return
	}
	defer db.Close()

	// testing DB connection with database
	if err := configs.TestDB(db); err != nil {
		log.Println("Ping to DB failsed\nCause: ", err.Error())
	}
	log.Println("DB connected")

	// inisialization redish
	rdb := configs.InitRedis()
	cmd := rdb.Ping(context.Background())
	if cmd.Err() != nil {
		log.Println("failed ping on redis \nCause:", cmd.Err().Error())
		return
	}
	log.Println("Redis Connected")
	defer rdb.Close()

	// Inisialization engine gin, HTTP framework
	router := routers.InitRouter(db, rdb)
	//  run the engine gin
	// Run this project on 127.0.0.1:8080 or localhost:8080
	if runtime.GOOS == "windows" {
		router.Run("127.0.0.1:8080")
	} else {
		router.Run(":8080")
	}
}
