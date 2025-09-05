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

// function add New users
func (a *AuthRepository) NewUser(reqContxt context.Context, email, password string) (models.User, error) {
	sql := "INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id, first_name, last_name, email, role, gender"
	values := []any{email, password}
	var NewUser models.User
	err := a.db.QueryRow(reqContxt, sql, values...).Scan(&NewUser.Id, &NewUser.FirstName, &NewUser.LastName, &NewUser.Email, &NewUser.Role, &NewUser.Gender)
	if err != nil {
		log.Println("Scan Error, ", err.Error())
		return models.User{}, err
	}
	return NewUser, nil
}
