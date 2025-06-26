package main

import (
	"context"
	"log"
	"time"

	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/config"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/adapters/inbound/http"
	outbound_mongo "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/adapters/outbound/mongo"
	outbound_security "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/adapters/outbound/security"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/application/service"
	"github.com/labstack/echo/v4"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.LoadMongoConfigFromEnv().URI))
	if err != nil {
		log.Fatal("Mongo connect error : ", err)
	}
	collection := client.Database("my_data").Collection("users")

	repo := outbound_mongo.NewUserMongoRepository(collection)
	hash := outbound_security.NewBcryptHasher()
	svc := service.NewUserService(repo, hash)
	handler := http.NewUserHandler(svc)

	e := echo.New()
	handler.RegisterRoutes(e)
	log.Fatal(e.Start(":8080"))
}
