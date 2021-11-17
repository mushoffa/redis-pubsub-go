package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mushoffa/go-library/redis"
)

type Publisher struct {
	Topic string      `json:"topic" binding:"required"`
	Data  interface{} `json:"data" binding:"required"`
}

func main() {
	redisHost := "localhost:6379"
	redisPassword := ""
	redisDB := 0

	redis, err := redis.NewRedisClient(redisHost, redisPassword, redisDB)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	fmt.Println("Redis PING: ", redis.Ping())

	pubsub := redis.PSubscribe("publisher.*")

	if _, err := pubsub.Receive(context.Background()); err != nil {
		log.Fatalf("Err: %v", err)
	}

	// defer pubsub.Close()

	channel := pubsub.Channel()

	go func() {
		for msg := range channel {

			// var payload Publisher

			// if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&payload); err != nil {
			// 	fmt.Println("Err: ", err)

			// 	continue
			// }

			fmt.Printf("Received message: %v, Data: %v", msg.Channel, msg.Payload)

		}
	}()

	for {
	}
}
