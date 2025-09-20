package repositories

import (
	"context"
	"errors"
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
	sql := `SELECT p.first_name, p.last_name, p.avatar_path, p.point, u.email, p.phone_number,  p.gender, p.created_at, p.updated_at
		FROM profiles p 
		JOIN users u ON p.user_id=u.id
		WHERE p.user_id =$1`

	var user models.User
	err := u.db.QueryRow(reqCntxt, sql, userid).Scan(&user.FirstName, &user.LastName, &user.AvatarPath, &user.Point, &user.Email, &user.Phone, &user.Gender, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Error when select, \nCause: ", err)
		return models.User{}, err
	}
	log.Println(user)
	return user, nil
}

// edit user, possible edit data : first_name, last_name, point, Gender, phone_number
func (u *UserRepository) EditUser(reqCntxt context.Context, body models.User, userID int) (models.User, error) {
	// query and validation which column will update
	values := []any{}
	sql := `UPDATE profiles SET `
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
	if body.Point != 0 {
		idx := strconv.Itoa(len(values) + 1)
		sql += "point=$" + idx + " ,"
		values = append(values, body.Point)
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

	idx := strconv.Itoa(len(values) + 1)
	sql += " updated_at=CURRENT_TIMESTAMP WHERE user_id =$" + idx + " RETURNING user_id, first_name, last_name, avatar_path, point, phone_number, gender, created_at, updated_at "
	values = append(values, userID)

	log.Println(values...)
	log.Println(sql)
	// execute the query into database, and bindit to new model user
	var newProfile models.User
	err := u.db.QueryRow(reqCntxt, sql, values...).Scan(&newProfile.Id, &newProfile.FirstName, &newProfile.LastName, &newProfile.AvatarPath, &newProfile.Point, &newProfile.Phone, &newProfile.Gender, &newProfile.CreatedAt, &newProfile.UpdatedAt)
	if err != nil {
		log.Println("scan Error. ", err.Error())
		return models.User{}, err
	}

	return newProfile, nil
}

// edit user avatar
func (u *UserRepository) EditAvatarProfile(reqCntxt context.Context, image string, id int) (models.User, error) {
	sql := "UPDATE profiles SET avatar_path=$1 , updated_at=CURRENT_TIMESTAMP WHERE user_id=$2 RETURNING user_id, first_name, last_name, avatar_path, point, phone_number, gender, created_at, updated_at"

	// execute the query into database, and bindit to new model user
	var newProfile models.User
	err := u.db.QueryRow(reqCntxt, sql, image, id).Scan(&newProfile.Id, &newProfile.FirstName, &newProfile.LastName, &newProfile.AvatarPath, &newProfile.Point, &newProfile.Phone, &newProfile.Gender, &newProfile.CreatedAt, &newProfile.UpdatedAt)
	if err != nil {
		log.Println("scan Error. ", err.Error())
		return models.User{}, err
	}

	return newProfile, nil
}

// function update password
// insert hashed password
// query effected : only table users effected
func (u *UserRepository) EditPasswordUser(reqCntxt context.Context, newHashedPass string, userId int) error {
	sql := `UPDATE users SET password=$1 WHERE id=$2`
	values := []any{newHashedPass, userId}

	cmd, err := u.db.Exec(reqCntxt, sql, values...)
	if err != nil {
		log.Println("Failed execute query update password users\nCause:", err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		log.Println("no row effected when insert movie_genres maybe failed?")
		return errors.New("no row effected when insert movie_genres maybe failed?")
	}

	// if no error, return error as nil
	return nil
}
