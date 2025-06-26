package inbound

import "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"

type UserService interface {
	CreateUser(user *domain.User) error
	GetUsers() ([]*domain.User, error)
}
