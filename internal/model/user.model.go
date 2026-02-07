package model

import "database/sql"

type User struct {
	ID     int            `db:"id"`
	Email  string         `db:"email"`
	Name   sql.NullString `db:"name"`
	Avatar sql.NullString `db:"avatar"`
	Bio    sql.NullString `db:"bio"`
}

type Users struct {
	ID   int            `db:"id"`
	Name sql.NullString `db:"name"`
}
