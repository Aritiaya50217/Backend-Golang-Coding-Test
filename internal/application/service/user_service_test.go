package service_test

import (
	"log"
	"testing"

	service "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/application/service"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockRepo struct{ mock.Mock }
type MockHasher struct{ mock.Mock }

func (m *MockRepo) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepo) GetUserById(id string) (*domain.User, error) {
	args := m.Called(id)
	user, _ := args.Get(0).(*domain.User)
	return user, args.Error(1)
}

func (m *MockRepo) GetUsers() ([]*domain.User, error) {
	args := m.Called()
	users, _ := args.Get(0).([]*domain.User)
	return users, args.Error(1)
}

func (m *MockRepo) Save(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepo) UpdateUser(id, name, email string) error {
	args := m.Called(id, name, email)
	return args.Error(0)
}

func (m *MockRepo) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) CountUsers() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) Compare(hashedPassword, password string) bool {
	args := m.Called(hashedPassword, password)
	return args.Bool(0)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockHasher := new(MockHasher)

	userService := &service.UserService{Repo: mockRepo, Hasher: mockHasher}

	user := &domain.User{
		Name:     "Test Create User",
		Email:    "test@gmail.com",
		Password: "test1111",
	}

	// ไม่มี user ซ้ำใน database
	mockRepo.On("GetUserByEmail", user.Email).Return(nil, nil)
	mockHasher.On("Hash", user.Password).Return("hashed pass", nil)
	mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil)

	err := userService.CreateUser(user)

	// Assertion
	assert.NoError(t, err)
	assert.Equal(t, "hashed pass", user.Password)

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockRepo)
	MockHasher := new(MockHasher)

	service := &service.UserService{Repo: mockRepo, Hasher: MockHasher}
	idString := "60b8d295f1a4e3e7d5a2b85f"
	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		log.Fatal(err)
	}
	mockRepo.On("GetUserByEmail", "test@gmail.com").Return(&domain.User{ID: id}, nil)

	err = service.CreateUser(&domain.User{
		Email: "test@gmail.com",
	})
	assert.Equal(t, err, "email already exists")
}
