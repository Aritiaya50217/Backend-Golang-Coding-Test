package mongo

import (
	"context"
	"fmt"
	"time"

	security "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/adapters/outbound/security"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"
	outbound "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/outbound"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongoRepository struct {
	col    *mongo.Collection
	hasher outbound.PasswordHasher
}

func NewUserMongoRepository(c *mongo.Collection) outbound.UserRepository {
	return &userMongoRepository{col: c}
}

func (r *userMongoRepository) InitDefaultUser() error {
	count, err := r.col.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return err
	}
	// hash
	hashed, err := security.NewBcryptHasher().Hash("admin1234")
	if err != nil {
		return err
	}

	if count == 0 {
		defaultUser := &domain.User{
			ID:        primitive.NewObjectID(),
			Name:      "Admin",
			Email:     "admin@example.com",
			Password:  hashed,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err = r.col.InsertOne(context.Background(), defaultUser)
		return err
	}
	return nil
}

func (r *userMongoRepository) CreateUser(user *domain.User) error {
	_, err := r.col.InsertOne(context.Background(), user)
	return err
}

func (r *userMongoRepository) GetUsers() ([]*domain.User, error) {
	cursor, err := r.col.Find(context.Background(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []*domain.User
	err = cursor.All(context.Background(), &users)
	return users, err
}

func (r *userMongoRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.col.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *userMongoRepository) GetUserById(id string) (*domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := r.col.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user); err != nil {
		return nil, err
	}

	user.ID = objID
	return &user, nil
}

func (r *userMongoRepository) UpdateUser(id, name, email string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":       name,
			"email":      email,
			"updated_at": now,
		},
	}
	_, err = r.col.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	return err
}

func (r *userMongoRepository) DeleteUser(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.col.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *userMongoRepository) CountUsers() (int64, error) {
	return r.col.CountDocuments(context.TODO(), bson.M{})
}
