package setting

import "fmt"

var MQ = struct {
	MessageType   string `ini:"MESSAGE_TYPE"`
	ServerAddr    string `ini:"SERVER_ADDR"`
	ServerVersion string `ini:"SERVER_VERSION"`
	SslPath       string `ini:"SSL_PATH"`
}{}

func loadMQFrom(rootCfg ConfigProvider) error {
	sec, _ := rootCfg.GetSection("message")
	if err := sec.MapTo(&MQ); err != nil {
		return fmt.Errorf("failed to map message queen settings: %v", err)
	}
	return nil
}
