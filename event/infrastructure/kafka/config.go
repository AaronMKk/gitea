package kafka

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"log"
	"os"
)

type Config struct {
	kfklib.Config
}

func SetDefault() Config {
	kafkaAddr := os.Getenv("KAFKA_ADDR")
	kafkaVer := os.Getenv("KAFKA_VER")

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
