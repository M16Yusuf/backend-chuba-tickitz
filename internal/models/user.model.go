package models

import (
	"mime/multipart"
	"time"
)

type User struct {
	Id         int        `db:"id" json:"user_id"`
	FirstName  string     `db:"first_name" json:"first_name"`
	LastName   string     `db:"last_name" json:"last_name"`
	AvatarPath *string    `db:"avatar_path" json:"profile_path"`
	Point      int        `db:"point" json:"point"`
	Phone      *string    `db:"phone_number" json:"phone"`
	Gender     *string    `db:"gender" json:"gender"`
	Email      string     `db:"email" json:"email,omitempty"`
	Password   string     `db:"password" json:"password,omitempty"`
	Role       string     `db:"role" json:"role,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
}

type UserAvatar struct {
	User
	Image *multipart.FileHeader `form:"avatar"`
}

// buat user untuk db dan response/body
// get profile/ history berdasarkan token di parsing buat dapat id user
// Update query
