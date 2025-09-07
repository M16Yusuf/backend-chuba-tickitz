package repositories

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// get data user
func (u *UserRepository) GetDataUser(reqCntxt context.Context, userid string) (models.User, error) {
	sql := `SELECT id, first_name, last_name, avatar_path, email, phone_number, password, point, gender
  	FROM users where id = $1 `

	var user models.User
	err := u.db.QueryRow(reqCntxt, sql, userid).Scan(&user.Id, &user.FirstName, &user.LastName, &user.AvatarPath, &user.Email, &user.Phone, &user.Password, &user.Point, &user.Gender)
	if err != nil {
		log.Println("Error when select, \nCause: ", err)
		return models.User{}, err
	}

	return user, nil
}
