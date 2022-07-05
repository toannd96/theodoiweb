package db

import (
	"context"
	"fmt"

	"analytics-api/configs"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func NewRedis() {
	if configs.IsDev() {
		configs.Redis.Client = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", configs.Redis.Host, configs.Redis.Port),
		})
	} else {
		opt, err := redis.ParseURL(configs.Redis.URL)
		if err != nil {
			logrus.Error(context.TODO(), err)
		}
		configs.Redis.Client = redis.NewClient(opt)
	}
	_, err := configs.Redis.Client.Ping().Result()
	if err != nil {
		logrus.Error(context.TODO(), err)
	}
}
