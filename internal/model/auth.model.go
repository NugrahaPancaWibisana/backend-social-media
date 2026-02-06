package model

type Account struct {
	ID          int           `db:"id"`
	Email       string        `db:"email"`
	Password    string        `db:"password"`
}
