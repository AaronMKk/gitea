package messagequeen

import (
	"fmt"

	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"
)

type MQConfig struct {
	MessageType   string `ini:"MESSAGE_TYPE"`
	ServerAddr    string `ini:"SERVER_ADDR"`
	ServerVersion string `ini:"SERVER_VERSION"`
	SslPath       string `ini:"SSL_PATH"`
}

type MessageType string

type NewMessageFunc func(cfg MQConfig) error

const (
	KafkaMessageType MessageType = "kafka"
)

var messageMap = map[MessageType]NewMessageFunc{}

// RegisterMessageType registers a provided storage type with a function to create it
func RegisterMessageType(typ MessageType, fn func(cfg MQConfig) error) {
	messageMap[typ] = fn
}

// Init the message queen, (ex: ActiveMQ、RocketMQ、RabbitMQ、Kafka)
func Init() (err error) {
	log.Info("Initialising message queen with type: %s", setting.MQ.MessageType)
	err = NewMessage(MessageType(setting.MQ.MessageType), setting.MQ)
	return err
}

// NewMessage takes a message queen type and some config and returns an error if exists
func NewMessage(typStr MessageType, cfg MQConfig) error {
	fn, ok := messageMap[typStr]
	if !ok {
		return fmt.Errorf("unsupported message queen type: %s", typStr)
	}

	err := fn(cfg)
	if err != nil {
		return err
	}
	return nil
}
