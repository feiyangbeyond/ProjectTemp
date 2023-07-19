package redis

import (
	"context"
	"fmt"

	"deviceback/v3/pkg/config"

	"github.com/go-redis/redis/v8"
)

func NewRdsClient(c config.RedisConfig) (*redis.Client, error) {
	rds := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", c.Host, c.Port),
		DB:           c.Db,
		DialTimeout:  0,
		ReadTimeout:  0,
		WriteTimeout: 0,
		Password:     c.Password,
	})

	_, err := rds.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return rds, nil
}
