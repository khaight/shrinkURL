package db

import (
	"log"

	"github.com/go-redis/redis"
)

// COUNTER global counter to increment for each url created
const COUNTER = "__counter__"

// Client redis client
type Client struct {
	*redis.Client
}

// Config redis client
type Config struct {
	Addr     string
	Password string
}

// InitConnection redis client
func InitConnection(config *Config) (*Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       0, //DEFAULT
	})

	client := &Client{cli}
	_, err := cli.Ping().Result()

	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	return client, nil
}

// IncrCounter Increments a global counter
func (client *Client) IncrCounter() (int64, error) {
	ctr, _ := client.Incr(COUNTER).Result()
	return ctr, nil
}
