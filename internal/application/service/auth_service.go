package service

import (
	"errors"

	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/inbound"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/outbound"
)

type authService struct {
	userRepo outbound.UserRepository
	hasher   outbound.PasswordHasher
	tokens   outbound.TokenGenerator
}

func NewAuthService(repo outbound.UserRepository, hasher outbound.PasswordHasher, tokens outbound.TokenGenerator) inbound.AuthenService {
	return &authService{
		userRepo: repo,
		hasher:   hasher,
		tokens:   tokens,
	}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid email or password")
	}
	if !s.hasher.Compare(user.Password, password) {
		return "", errors.New("invalid email or password")
	}

	// genarate token
	token, err := s.tokens.GenerateToken(user.ID.Hex())
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *authService) Authorize(userID, action string) (bool, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	switch action {
	case "create_user":
		return user.Role == domain.RoleAdmin, nil
	case "register":
		return true, nil
	default:
		return false, nil
	}
}

func (s *authService) CreateUser(email, password string) error {
	user := &domain.User{
		Email:    email,
		Password: password,
		Role:     domain.RoleUser,
	}
	return s.userRepo.CreateUser(user)
}
