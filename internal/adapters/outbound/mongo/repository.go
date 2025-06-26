package mongo

import (
	"context"

	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"
	outbound "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/outbound"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongoRepository struct {
	col *mongo.Collection
}

func NewUserMongoRepository(c *mongo.Collection) outbound.UserRepository {
	return &userMongoRepository{col: c}
}

func (r *userMongoRepository) Save(user *domain.User) error {
	_, err := r.col.InsertOne(context.Background(), user)
	return err
}

func (r *userMongoRepository) FindAll() ([]*domain.User, error) {
	cursor, err := r.col.Find(context.Background(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []*domain.User
	err = cursor.All(context.Background(), &users)
	return users, err
}

func (r *userMongoRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.col.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}
