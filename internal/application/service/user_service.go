package service

import (
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/inbound"
	outbound "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/outbound"
)

type userService struct {
	repo outbound.UserRepository
}

func NewUserService(r outbound.UserRepository) inbound.UserService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(user *domain.User) error {
	return s.repo.Save(user)
}

func (s *userService) GetUsers() ([]*domain.User, error) {
	return s.repo.FindAll()
}
