package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/utils"
	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAuthRepository(db *pgxpool.Pool, rdb *redis.Client) *AuthRepository {
	return &AuthRepository{db: db, rdb: rdb}
}

// function to get user data
func (a *AuthRepository) GetUserWithEmail(reqContxt context.Context, email string) (models.User, error) {
	// get user by input email and validate user
	// sql := "SELECT id, email, password, role FROM users WHERE email=$1"
	sql := `SELECT u.id, u.email, u.password, u.role, p.first_name, p.last_name, p.avatar_path
		FROM users u 
		JOIN profiles p ON p.user_id = u.id
		WHERE email=$1`

	var User models.User
	if err := a.db.QueryRow(reqContxt, sql, email).Scan(&User.Id, &User.Email, &User.Password, &User.Role, &User.FirstName, &User.LastName, &User.AvatarPath); err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		return models.User{}, err
	}
	return User, nil
}

// function add New users
// inputs : Validated and hashed email and password
// Query tables effected : users, profiles
func (a *AuthRepository) NewUser(reqContxt context.Context, email, password string) (models.User, error) {

	// insert all query inside postgreSQL's transaction
	tx, err := a.db.Begin(reqContxt)
	if err != nil {
		log.Println("Failed to begin DB transaction\nCause: ", err)
		return models.User{}, err
	}
	defer tx.Rollback(reqContxt)

	// insert inputs new user to table user
	sql1 := "INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id"
	values := []any{email, password}
	var tempNewUserID models.User
	if err := tx.QueryRow(reqContxt, sql1, values...).Scan(&tempNewUserID.Id); err != nil {
		log.Println("Scan Error, ", err.Error())
		return models.User{}, err
	}

	// insert new user_id as new profile
	sql2 := "INSERT INTO profiles(user_id) VALUES($1) RETURNING user_id, first_name, last_name, avatar_path, point, phone_number, gender"
	var NewUser models.User
	if err := tx.QueryRow(reqContxt, sql2, tempNewUserID.Id).Scan(&NewUser.Id, &NewUser.FirstName, &NewUser.LastName, &NewUser.AvatarPath, &NewUser.Point, &NewUser.Phone, &NewUser.Gender); err != nil {
		log.Println("Scan Error, ", err.Error())
		return models.User{}, err
	}

	// commit transaction if both query success execute
	if err := tx.Commit(reqContxt); err != nil {
		log.Println("Failed to commit DBtransaction\nCause: ", err)
		return models.User{}, err
	}
	log.Println("success to commit DB transaction")

	// return result returning from query as model user
	return NewUser, nil
}

// function repo to redis db
// blacklist token user
func (a *AuthRepository) BlaclistToken(reqContxt context.Context, token string) error {
	// use utils.BlacklistToken for logout token
	if err := utils.BlaclistTokenRedish(reqContxt, *a.rdb, token); err != nil {
		log.Println("failed blacklist token, ", err)
		return err
	}
	// is success return nil
	return nil
}
