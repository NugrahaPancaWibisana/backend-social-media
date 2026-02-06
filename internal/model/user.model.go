package model

import "database/sql"

type User struct {
	ID          int           `db:"id"`
	Email       string        `db:"email"`
	Password    string        `db:"password"`
	LastLoginAt *sql.NullTime `db:"lastlogin_at"`
}
