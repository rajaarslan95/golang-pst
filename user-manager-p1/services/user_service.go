package services

import (
	"user-manager/schemas"
	"user-manager/store"
)

type UserService struct {
	UserStore store.UserStore
}

func (s *UserService) Create(user schemas.User) error {
	return s.UserStore.AddUser(user)
}

func (s *UserService) Update(user schemas.User) error {
	_, err := s.UserStore.GetUser(user.ID)
	if err != nil {
		return err
	}
	return s.UserStore.UpdateUser(user)
}

func (s *UserService) Get(id int) (schemas.User, error) {
	return s.UserStore.GetUser(id)
}

func (s *UserService) Delete(id int) error {
	_, err := s.UserStore.GetUser(id)
	if err != nil {
		return err
	}
	return s.UserStore.DeleteUser(id)
}
