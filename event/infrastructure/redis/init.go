package redis

import (
	redislib "github.com/opensourceways/redis-lib"
)

func Init(cfg *Config) error {
	return redislib.Init(&cfg.Config, false)

}

func Close() error {
	return redislib.Close()
}
