package repositories

import (
	"context"
	"log"
	"strconv"

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
func (u *UserRepository) GetDataUser(reqCntxt context.Context, userid int) (models.User, error) {
	sql := `SELECT id, first_name, last_name, avatar_path, email, phone_number, point, gender, updated_at
  	FROM users where id = $1 `

	var user models.User
	err := u.db.QueryRow(reqCntxt, sql, userid).Scan(&user.Id, &user.FirstName, &user.LastName, &user.AvatarPath, &user.Email, &user.Phone, &user.Point, &user.Gender, &user.UpdatedAt)
	if err != nil {
		log.Println("Error when select, \nCause: ", err)
		return models.User{}, err
	}

	return user, nil
}

// edit user
func (u *UserRepository) EditUser(reqCntxt context.Context, body models.User, userID int) (models.User, error) {
	// query and validation which column will update
	values := []any{}
	sql := `UPDATE users SET `
	if body.FirstName != "" {
		idx := strconv.Itoa(len(values) + 1)
		sql += "first_name=$" + idx + " ,"
		values = append(values, body.FirstName)
	}
	if body.LastName != "" {
		idx := strconv.Itoa(len(values) + 1)
		sql += "last_name=$" + idx + " ,"
		values = append(values, body.LastName)
	}
	if body.Email != "" {
		idx := strconv.Itoa(len(values) + 1)
		sql += "email=$" + idx + " ,"
		values = append(values, body.Email)
	}
	if body.Gender != nil && *body.Gender != "" {
		idx := strconv.Itoa(len(values) + 1)
		sql += "gender=$" + idx + " ,"
		values = append(values, body.Gender)
	}
	if body.Phone != nil {
		idx := strconv.Itoa(len(values) + 1)
		sql += "phone_number=$" + idx + " ,"
		values = append(values, body.Phone)
	}
	if body.AvatarPath != nil {
		idx := strconv.Itoa(len(values) + 1)
		sql += "avatar_path=$" + idx + " ,"
		values = append(values, body.AvatarPath)
	}
	if body.Point != 0 {
		idx := strconv.Itoa(len(values) + 1)
		sql += "point=$" + idx + " ,"
		values = append(values, body.Point)
	}
	if body.Role != "" {
		idx := strconv.Itoa(len(values) + 1)
		sql += "role=$" + idx + " ,"
		values = append(values, body.Role)
	}

	if body.Password != "" {
		idx := strconv.Itoa(len(values) + 1)
		sql += "password=$" + idx + " ,"
		values = append(values, body.Password)
	}
	idx := strconv.Itoa(len(values) + 1)
	sql += " updated_at=CURRENT_TIMESTAMP WHERE id =$" + idx + " RETURNING id, first_name, last_name, avatar_path, email, point, role, phone_number, password, gender, updated_at "
	values = append(values, userID)

	log.Println(values...)
	log.Println(sql)
	// execute the query into database, and bindit to new model user
	var newProfile models.User
	err := u.db.QueryRow(reqCntxt, sql, values...).Scan(&newProfile.Id, &newProfile.FirstName, &newProfile.LastName, &newProfile.AvatarPath, &newProfile.Email, &newProfile.Point, &newProfile.Role, &newProfile.Phone, &newProfile.Password, &newProfile.Gender, &newProfile.UpdatedAt)
	if err != nil {
		log.Println("scan Error. ", err.Error())
		return models.User{}, err
	}

	return newProfile, nil
}
