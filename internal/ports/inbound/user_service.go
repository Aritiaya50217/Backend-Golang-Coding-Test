package inbound

import "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"

type UserService interface {
	CreateUser(user *domain.User) error
	GetUser(id string) (*domain.User, error)
	GetUsers() ([]*domain.User, error)
	UpdateUser(id, name, email string) error
}
