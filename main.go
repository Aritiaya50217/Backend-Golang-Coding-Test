package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/config"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/adapters/inbound/http"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/adapters/inbound/http/middleware"
	outbound_mongo "github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/adapters/outbound/mongo"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/adapters/outbound/security"
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

	// secret key
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	userRepo := outbound_mongo.NewUserMongoRepository(collection)
	hash := outbound_security.NewBcryptHasher()
	uservice := service.NewUserService(userRepo, hash)
	userHandler := http.NewUserHandler(uservice)

	// auth
	tokenGen := security.NewJWTToKenGenarator(secret)
	authService := service.NewAuthService(userRepo, hash, tokenGen)
	authHandler := http.NewAuthHandler(authService)

	e := echo.New()
	// login
	authHandler.RegisterRoutes(e)

	// middleware
	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware)
	api.Use(middleware.LoggingMiddleware)
	// router
	userHandler.RegisterRoutes(api)
	log.Fatal(e.Start(":8080"))
}
