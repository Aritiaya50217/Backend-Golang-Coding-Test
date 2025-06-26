package service

import (
	"context"
	"log"
	"time"

	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/outbound"
)

func StartUserCountLogger(ctx context.Context, repo outbound.UserRepository) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				count, err := repo.CountUsers()
				if err != nil {
					log.Printf("[Error] count users :  %v", err)
					continue
				}
				log.Printf("[INFO] total users: %d", count)
			case <-ctx.Done():
				log.Println("[INFO] stop user count logger")
				ticker.Stop()
				return
			}
		}
	}()
}
