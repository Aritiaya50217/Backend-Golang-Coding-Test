package outbound

import (
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"
)

type UserRepository interface {
	Save(user *domain.User) error
	GetUserById(id string) (*domain.User, error)
	GetUsers() ([]*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(id, name, email string) error
	DeleteUser(id string) error
	CountUsers() (int64, error)
}
