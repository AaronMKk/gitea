package kafka

import (
	"os"

	kfklib "github.com/opensourceways/kafka-lib/agent"
)

type Config struct {
	kfklib.Config
}

func SetDefault() Config {
	return Config{
		kfklib.Config{
			Address:        os.Getenv("KAFKA_ADDR"),
			Version:        os.Getenv("KAFKA_VER"),
			SkipCertVerify: true,
		},
	}
}
