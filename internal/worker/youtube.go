package worker

import (
	initializers "github.com/souvik150/Fampay-Backend-Assignment/config"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/services"
	"log"
	"time"
)

func StartVideoWorker(topic string) {
	config, _ := initializers.LoadConfig(".")
	apiKeys := config.APIKeysYT

	keyIndex := 0

	go func() {
		for {
			apiKey := apiKeys[keyIndex]

			videos, err := services.CallYouTubeAPI(apiKey, topic)
			if err != nil {
				log.Printf("Error fetching videos: %v", err)
				// Switch to the next available API key if quota is exhausted
				keyIndex = (keyIndex + 1) % len(apiKeys)
				continue
			}

			for i := range videos {
				err := services.SaveVideo(videos[i])
				if err != nil {
					log.Println("Error storing videos:", err)
				}
			}

			time.Sleep(30 * time.Second)
		}
	}()
}
