package service

import (
	"context"
	"errors"
	"time"

	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/adapters/outbound/security"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/inbound"
	outbound "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/outbound"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	Repo   outbound.UserRepository
	Hasher outbound.PasswordHasher
}

func NewUserService(r outbound.UserRepository, h outbound.PasswordHasher) inbound.UserService {
	return &UserService{Repo: r, Hasher: h}
}

func (s *UserService) CreateUser(user *domain.User) error {
	// check email
	existing, _ := s.Repo.GetUserByEmail(user.Email)
	if existing != nil {
		return errors.New("email already exists")
	}
	// hash
	hashed, err := s.Hasher.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// save to repository
	return s.Repo.Save(user)
}

func (s *UserService) GetUser(id string) (*domain.User, error) {
	return s.Repo.GetUserById(id)
}

func (s *UserService) GetUsers() ([]*domain.User, error) {
	return s.Repo.GetUsers()
}

func (s *UserService) UpdateUser(id, name, email string) error {
	user, err := s.Repo.GetUserById(id)
	if err != nil {
		return errors.New("user not found")
	}
	existing, _ := s.Repo.GetUserByEmail(email)
	if existing != nil && existing.ID != user.ID {
		return errors.New("email already is use")
	}
	return s.Repo.UpdateUser(id, name, email)
}

func (s *UserService) DeleteUser(id string) error {
	return s.Repo.DeleteUser(id)
}

func (s *UserService) CountUsers() error {
	_, err := s.Repo.CountUsers()
	return err
}

func (s *UserService) InitDefaultUser(ctx context.Context) error {
	existing, err := s.Repo.GetUserByEmail("admin@example.com")
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	// hash
	hashed, err := security.NewBcryptHasher().Hash("admin1234")
	if err != nil {
		return err
	}

	defaultUser := &domain.User{
		ID:        primitive.NewObjectID(),
		Name:      "Admin",
		Email:     "admin@example.com",
		Password:  hashed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.Repo.Save(defaultUser)
}
