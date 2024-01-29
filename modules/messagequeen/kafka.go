package messagequeen

import (
	"github.com/sirupsen/logrus"

	"code.gitea.io/gitea/modules/log"
	kfklib "github.com/opensourceways/kafka-lib/agent"
)

const (
	queueName = "gitea-kafka-queue"
)

type Config struct {
	kfklib.Config
}

func init() {
	RegisterMessageType(KafkaMessageType, newKafkaMessageQueue)
}

func retriveConfig(cfg MQConfig) Config {
	kafkaAddr := cfg.ServerAddr
	kafkaVer := cfg.ServerVersion

	// Check if KAFKA_VER is set in the environment
	if kafkaVer == "" {
		log.Fatal("KAFKA_VER is not set in the environment. " +
			"It's crucial to set this to avoid protocol version mismatches and ensure backward compatibility.")
	}

	return Config{
		kfklib.Config{
			Address:        kafkaAddr,
			Version:        kafkaVer,
			SkipCertVerify: true,
		},
	}
}

// newKafkaMessageQueue sets up a new Kafka message queue
func newKafkaMessageQueue(cfg MQConfig) error {
	var localConfig = retriveConfig(cfg)
	mqLog := logrus.NewEntry(logrus.StandardLogger())
	err := kfklib.Init(&localConfig.Config, mqLog, nil, queueName, true)
	if err != nil {
		return err
	}
	return nil
}
