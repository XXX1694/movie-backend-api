package worker

import (
	"log"
	"time"

	"golang/internal/repository"
)

func StartUserCountWorker(repos *repository.Repositories) {
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			users, err := repos.GetUsers(1000, 0)
			if err != nil {
				log.Printf("[Worker] Error counting users: %v", err)
				continue
			}
			log.Printf("[Worker] Total users in DB: %d", len(users))
		}
	}()
	log.Println("[Worker] User count worker started")
}
