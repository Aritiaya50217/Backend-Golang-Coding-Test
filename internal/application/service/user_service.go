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
	user.UpdatedAt = time.Now()

	// save to repository
	return s.repo.Save(user)
}

func (s *userService) GetUser(id string) (*domain.User, error) {
	return s.repo.GetUserById(id)
}

func (s *userService) GetUsers() ([]*domain.User, error) {
	return s.repo.GetUsers()
}

func (s *userService) UpdateUser(id, name, email string) error {
	user, err := s.repo.GetUserById(id)
	if err != nil {
		return errors.New("user not found")
	}
	existing, _ := s.repo.GetUserByEmail(email)
	if existing != nil && existing.ID != user.ID {
		return errors.New("email already is use")
	}
	return s.repo.UpdateUser(id, name, email)
}
