package models

import "database/sql"

type UserRole struct {
	Id int `db:"id"`
	RoleId int `db:"role_id"`
	UserCode sql.NullString `db:"user_code"`
}
