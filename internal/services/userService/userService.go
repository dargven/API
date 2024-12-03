package userService

import "API/internal/Storage/postrgeSQL"

type UserService struct {
	DB *postrgeSQL.Database
}

func (s *UserService) isEmailUnique(email string) (bool, error) {
	return true, nil
}

//Доделать
