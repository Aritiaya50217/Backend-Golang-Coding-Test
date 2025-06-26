package service

import (
	"errors"
	"time"

	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/inbound"
	outbound "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/outbound"
)

type userService struct {
	repo   outbound.UserRepository
	hasher outbound.PasswordHasher
}

func NewUserService(r outbound.UserRepository, h outbound.PasswordHasher) inbound.UserService {
	return &userService{repo: r, hasher: h}
}

func (s *userService) CreateUser(user *domain.User) error {
	// check email
	existing, _ := s.repo.GetUserByEmail(user.Email)
	if existing != nil {
		return errors.New("email already exists")
	}
	// hash
	hashed, err := s.hasher.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	user.CreatedAt = time.Now()
	// save to repository
	return s.repo.Save(user)
}

func (s *userService) GetUser(id string) (*domain.User, error) {
	return s.repo.GetUserById(id)
}

func (s *userService) GetUsers() ([]*domain.User, error) {
	return s.repo.GetUsers()
}
