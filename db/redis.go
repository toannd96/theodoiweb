package db

import (
	"analytics-api/configs"
	"analytics-api/internal/pkg/log"
	"context"

	"github.com/go-redis/redis"
)

func NewRedis() {
	// local
	// configs.Redis.Client = redis.NewClient(&redis.Options{
	// 	Addr: fmt.Sprintf("%s:%s", configs.Redis.Host, configs.Redis.Port),
	// })

	// heroku
	opt, err := redis.ParseURL(configs.Redis.URL)
	if err != nil {
		log.LogError(context.TODO(), err)
	}
	configs.Redis.Client = redis.NewClient(opt)

	_, err = configs.Redis.Client.Ping().Result()
	if err != nil {
		log.LogError(context.TODO(), err)
	}
}
