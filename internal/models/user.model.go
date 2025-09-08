package models

import "time"

type User struct {
	Id         int        `db:"id" json:"user_id"`
	FirstName  string     `db:"first_name" json:"first_name"`
	LastName   string     `db:"last_name" json:"last_name"`
	AvatarPath string     `db:"avatar_path" json:"profile_path,omitempty"`
	Email      string     `db:"email" json:"email"`
	Point      int        `db:"point" json:"point,omitempty"`
	Role       string     `db:"role" json:"role"`
	Phone      string     `db:"phone_number" json:"phone,omitempty"`
	Password   string     `db:"password" json:"password,omitempty"`
	Gender     *string    `db:"gender" json:"gender,omitempty"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
