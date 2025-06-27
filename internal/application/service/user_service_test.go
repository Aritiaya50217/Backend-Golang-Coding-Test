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
type MockCollection struct{ mock.Mock }

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

func BenchmarkCreateUser(b *testing.B) {
	mockRepo := new(MockRepo)
	mockHasher := new(MockHasher)
	userService := &service.UserService{Repo: mockRepo, Hasher: mockHasher}

	mockRepo.On("GetUserByEmail", mock.Anything).Return(nil, nil)
	mockHasher.On("Hash", mock.Anything).Return("hashed-pass", nil)
	mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil)

	for i := 0; i < b.N; i++ {
		user := &domain.User{
			Name:     "Test Create User",
			Email:    "test@gmail.com",
			Password: "test1111",
		}
		_ = userService.CreateUser(user)
	}

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

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockHasher := new(MockHasher)

	userService := &service.UserService{Repo: mockRepo, Hasher: mockHasher}

	userID := "60b8d295f1a4e3e7d5a2b85f"
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		t.Fatalf("invalid ObjectID hex: %v", err)
	}

	mockRepo.On("GetUserById", mock.AnythingOfType("string")).Return(&domain.User{
		ID:    id,
		Name:  "John Doe",
		Email: "john@example.com",
	}, nil)

	user, err := userService.GetUser(userID)

	assert.NoError(t, err)
	assert.Equal(t, id, user.ID)

	mockRepo.AssertExpectations(t)
}

func BenchmarkGetUserById(b *testing.B) {
	mockRepo := new(MockRepo)
	userService := &service.UserService{Repo: mockRepo}
	userID := "60b8d295f1a4e3e7d5a2b85f"
	id, _ := primitive.ObjectIDFromHex(userID)
	user := &domain.User{
		ID: id,
	}

	mockRepo.On("GetUserById", userID).Return(user, nil)

	for i := 0; i < b.N; i++ {
		_, _ = userService.GetUser(userID)
	}
}

func TestGetUsers_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	service := &service.UserService{Repo: mockRepo}

	userFirst := "60b8d295f1a4e3e7d5a2b85f"
	userSecond := "685d9aaaabac97c7040575f1"
	userId1, err := primitive.ObjectIDFromHex(userFirst)
	if err != nil {
		t.Fatalf("invalid ObjectID hex: %v", err)
	}

	userId2, err := primitive.ObjectIDFromHex(userSecond)
	if err != nil {
		t.Fatalf("invalid ObjectID hex: %v", err)
	}

	users := []*domain.User{{ID: userId1}, {ID: userId2}}
	mockRepo.On("GetUsers").Return(users, nil)

	result, err := service.GetUsers()
	assert.NoError(t, err)
	assert.Len(t, result, 2)

}

func BenchmarkGetUsers(b *testing.B) {
	mockRepo := new(MockRepo)
	userService := &service.UserService{Repo: mockRepo}

	userFirst := "60b8d295f1a4e3e7d5a2b85f"
	userSecond := "685d9aaaabac97c7040575f1"
	userId1, _ := primitive.ObjectIDFromHex(userFirst)
	userId2, _ := primitive.ObjectIDFromHex(userSecond)

	users := []*domain.User{{ID: userId1}, {ID: userId2}}
	mockRepo.On("GetUsers").Return(users, nil)

	for i := 0; i < b.N; i++ {
		_, _ = userService.GetUsers()
	}
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	service := &service.UserService{Repo: mockRepo}

	userID := "60b8d295f1a4e3e7d5a2b85f"
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		t.Fatalf("invalid ObjectID hex: %v", err)
	}
	mockUser := &domain.User{
		ID:    id,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	mockRepo.On("GetUserById", userID).Return(mockUser, nil)
	mockRepo.On("GetUserByEmail", "new@example.com").Return(nil, nil)
	mockRepo.On("UpdateUser", userID, "New Name", "new@example.com").Return(nil)

	err = service.UpdateUser(userID, "New Name", "new@example.com")
	assert.NoError(t, err)
}

func BenchmarkUpdateUser(b *testing.B) {
	mockRepo := new(MockRepo)
	userService := &service.UserService{Repo: mockRepo}

	userID := "60b8d295f1a4e3e7d5a2b85f"
	id, _ := primitive.ObjectIDFromHex(userID)

	mockUser := &domain.User{
		ID:    id,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	mockRepo.On("GetUserById", userID).Return(mockUser, nil)
	mockRepo.On("GetUserByEmail", mock.Anything).Return(nil, nil)
	mockRepo.On("UpdateUser", userID, mock.Anything, mock.Anything).Return(nil)

	for i := 0; i < b.N; i++ {
		_ = userService.UpdateUser(userID, "New Name", "new@example.com")
	}
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	service := &service.UserService{Repo: mockRepo}
	userID := "60b8d295f1a4e3e7d5a2b85f"
	mockRepo.On("DeleteUser", userID).Return(nil)

	err := service.DeleteUser(userID)
	assert.NoError(t, err)
}

func BenchmarkDeleteUser(b *testing.B) {
	mockRepo := new(MockRepo)
	userService := &service.UserService{Repo: mockRepo}

	userID := "60b8d295f1a4e3e7d5a2b85f"
	mockRepo.On("DeleteUser", userID).Return(nil)

	for i := 0; i < b.N; i++ {
		_ = userService.DeleteUser(userID)
	}
}
