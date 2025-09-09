package models

import "time"

type User struct {
	Id         int        `db:"id" json:"user_id"`
	FirstName  string     `db:"first_name" json:"first_name"`
	LastName   string     `db:"last_name" json:"last_name"`
	AvatarPath *string    `db:"avatar_path" json:"profile_path"`
	Email      string     `db:"email" json:"email"`
	Point      int        `db:"point" json:"point"`
	Role       string     `db:"role" json:"role,omitempty"`
	Phone      *string    `db:"phone_number" json:"phone,omitempty"`
	Password   string     `db:"password" json:"password,omitempty"`
	Gender     *string    `db:"gender" json:"gender"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// buat user untuk db dan response/body
// get profile/ history berdasarkan token di parsing buat dapat id user
// Update query
