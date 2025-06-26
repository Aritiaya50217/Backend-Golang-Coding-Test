package config

import (
	"os"
	"strconv"
	"time"
)

type MongoConfig struct {
	URI     string
	Timeout time.Duration
}

func LoadMongoConfigFromEnv() MongoConfig {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	timeout := os.Getenv("MONGO_TIMEOUT_SEC")
	timeoutSec, err := strconv.Atoi(timeout)
	if err != nil || timeoutSec <= 0 {
		timeoutSec = 10
	}
	return MongoConfig{
		URI:     uri,
		Timeout: time.Duration(timeoutSec) * time.Second,
	}
}
