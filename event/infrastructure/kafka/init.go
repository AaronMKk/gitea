package kafka

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/kafka-lib/mq"
)

const (
	queueName = "gitea-kafka-queue"
)

func Init(cfg *Config, log mq.Logger, redis kfklib.Redis) error {
	return kfklib.Init(&cfg.Config, log, redis, queueName, true)
}

func Exit() {
	kfklib.Exit()
}
