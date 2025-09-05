package models

type User struct {
	Id         int    `db:"id" json:"user_id"`
	FirstName  string `db:"first_name" json:"first_name"`
	LastName   string `db:"last_name" json:"last_name"`
	AvatarPath string `db:"avatar_path" json:"profile_path"`
	Email      string `db:"email" json:"email"`
	Point      int    `db:"point" json:"point"`
	Role       string `db:"role" json:"role"`
	Phone      string `db:"phone_number" json:"phone"`
	Password   string `db:"password" json:"password"`
	Gender     string `db:"gender" json:"gender"`
}
