package configs

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() (*pgxpool.Pool, error) {
	// data database berada di .env
	dbUser := os.Getenv("DB_USER")
	dbUserPass := os.Getenv("DB_USER_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// inisilisasi DB
	connetString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbUserPass, dbHost, dbPort, dbName)
	return pgxpool.New(context.Background(), connetString)
}

func TestDB(db *pgxpool.Pool) error {
	return db.Ping(context.Background())
}
