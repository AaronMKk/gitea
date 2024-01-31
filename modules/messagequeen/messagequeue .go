package messagequeen

import (
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"
)

type MQConfig struct {
	MessageType   string `ini:"MESSAGE_TYPE"`
	ServerAddr    string `ini:"SERVER_ADDR"`
	ServerVersion string `ini:"SERVER_VERSION"`
	SslPath       string `ini:"SSL_PATH"`
	TopicName     string `ini:"TOPIC_NAME"`
}

type MessageType string

type NewMessageFunc func(cfg MQConfig) error

// Init the message queen, (ex: ActiveMQ、RocketMQ、RabbitMQ、Kafka)
func Init() (err error) {
	log.Info("Initialising message queen with type: %s", setting.MQ.MessageType)
	err = NewMessage(MQConfig(setting.MQ))
	return err
}

// NewMessage takes a message queen type and some config and returns an error if exists
func NewMessage(cfg MQConfig) error {
	return newKafkaMessageQueue(cfg)
}
