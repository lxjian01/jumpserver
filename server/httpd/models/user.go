package models

import "database/sql"

type User struct {
	Id int `db:"id"`
	UserCode string   `db:"user_code"`
	UserName sql.NullString `db:"user_name"`
}
