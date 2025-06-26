package outbound

import "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"

type UserRepository interface {
	Save(user *domain.User) error
	FindAll() ([]*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
}
