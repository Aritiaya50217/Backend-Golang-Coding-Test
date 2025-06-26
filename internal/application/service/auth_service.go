package service

import (
	"errors"

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
	user, err := s.userRepo.FindByEmail(email)
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
