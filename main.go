package main

import (
	"context"
	"log"
	"time"

	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/config"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.LoadMongoConfigFromEnv().URI))
	if err != nil {
		log.Fatal("Mongo connect error : ", err)
	}
	collection := client.Database("golangtest").Collection("users")

	_ = collection

	// repo := mongo.NewUserMongoRepository(collection)
	// svc := service.NewUserService(repo)
	// handler := http.NewUserHandler(svc)

	e := echo.New()
	// handler.RegisterRoutes(e)
	log.Fatal(e.Start(":8080"))
}
