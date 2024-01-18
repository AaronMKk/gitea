package redis

import (
	"time"

	redislib "github.com/opensourceways/redis-lib"
)

const (
	defaultAddress  = "merlin-server-redis-1:6379"
	defaultPassword = "merlin"
	defaultDB       = 5
	defaultTimeOut  = 100 * time.Second
)

type Config struct {
	redislib.Config
}

func SetDefault() Config {
	return Config{
		redislib.Config{
			Address:  defaultAddress,
			Password: defaultPassword,
			DB:       defaultDB,
			Timeout:  int64(defaultTimeOut),
		},
	}
}
