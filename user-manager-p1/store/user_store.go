package store

import (
	"user-manager/schemas"
)

type UserStore interface {
	AddUser(user schemas.User) error
	UpdateUser(user schemas.User) error
	DeleteUser(id int) error
	GetUser(id int) (schemas.User, error)
}
