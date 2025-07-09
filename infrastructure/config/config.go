package config
import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	AccessToken string	`env:"MCP_ACCESS_TOKEN"`
	// Shopware API configuration
	TargetDomain string `env:"SHOPWARE_PLATFORM_DOMAIN"`
	ClientAccessKeyID   string `env:"SHOPWARE_ACCESS_KEY_ID"`
	ClientSecret	 string `env:"SHOPWARE_SECRET_ACCESS_KEY"`
}

func NewFromEnv() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
