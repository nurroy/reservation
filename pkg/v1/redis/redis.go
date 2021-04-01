package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

// Config is the struct to pass the Postgres configuration.
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Key      string
}

const REDIS_DB = 0

type ConfigStore interface {
	GetKey(string) string
}

type Token struct {
	Value  string
	Expire time.Duration
}

var client *redis.Client

// New initializes the redis
func New(config ConfigStore) error {
	var redisConfig RedisConfig

	// set redis config from env
	redisConfig = RedisConfig{
		Host:     config.GetKey("redis.host"),
		Port:     config.GetKey("redis.port"),
		Password: config.GetKey("redis.password"),
		Key:      config.GetKey("redis.key"),
	}

	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: fmt.Sprintf("%s", redisConfig.Password), // no password set
		DB:       REDIS_DB,                                // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}