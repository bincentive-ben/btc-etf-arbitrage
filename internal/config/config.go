package config

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type IBKRConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	IsPaper  bool   `mapstructure:"isPaper"`
}

type Config struct {
	IBKR IBKRConfig `mapstructure:"IBKR"`
}

type IBKRClient struct {
	client  *resty.Client
	BaseUrl string
}

var AppConfig Config

// LoadConfig reads the configuration from a YAML file
func LoadConfig(configFile string) error {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal the config into the AppConfig variable
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("error unmarshalling config: %v", err)
	}
	return nil
}

func GetConfig() Config {
	return AppConfig
}
