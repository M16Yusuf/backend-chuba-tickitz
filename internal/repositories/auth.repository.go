package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

// function to get user data
func (a *AuthRepository) GetUserWithEmail(reqContxt context.Context, email string) (models.User, error) {
	// get user by input email and validate user
	sql := "SELECT id, email, password, role FROM users WHERE email=$1"

	var User models.User
	if err := a.db.QueryRow(reqContxt, sql, email).Scan(&User.Id, &User.Email, &User.Password, &User.Role); err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		return models.User{}, err
	}
	return User, nil
}
