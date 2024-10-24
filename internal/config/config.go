package config

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var ConfigFile string
var TraderConfigFile string

var appConfig AppConfig
var traderConfig TraderConfig

type AppConfig struct {
	ReloadTraderConfigSeconds int           `mapstructure:"ReloadTraderConfigSeconds"`
	IbkrConfig                IbkrConfig    `mapstructure:"IbkrConfig"`
	BinanceConfig             BinanceConfig `mapstructure:"BinanceCOnfig"`
}

type DiscordConfig struct {
	Enable     bool   `mapstructure:"Enable"`
	WebhookUrl string `mapstructure:"WebhookUrl"`
	Intervarl  int    `mapstructure:"Interval"`
}
type IbkrConfig struct {
	RestEndpoint   string   `mapstructure:"RestEndpoint"`
	WsEndpoint     string   `mapstructure:"WsEndpoint"`
	Fields         []string `mapstructure:"Fields"`
	ContractIDList []string `mapstructure:"ContractIDList"`
	TickerList     []string `mapstructure:"TickerList"`
}

type BinanceConfig struct {
	ExchangeID   string `mapstructure:"ExchangeID"`
	RestEndpoint string `mapstructure:"RestEndpoint"`
	WsEndpoint   string `mapstructure:"WsEndpoint"`
	ApiKey       string `mapstructure:"ApiKey"`
	SecretKey    string `mapstructure:"SecretKey"`
	Type         string `mapstructure:"Type"`
}

type TraderConfig struct {
	ProfitRate   string  `mapstructure:"ProfitRate"`
	UsdtUsdRate  float64 `mapstructure:"USDT/USD"`
	IbitBtcRatio float64 `mapstructure:"IBIT/BTC"`
}

type IbkrClient struct {
	client  *resty.Client
	BaseUrl string
}

// LoadAppConfig reads the configuration from a YAML file
func LoadAppConfig(configFile string) error {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal the config into the appConfig variable
	if err := viper.Unmarshal(&appConfig); err != nil {
		return fmt.Errorf("error unmarshalling config: %v", err)
	}
	return nil
}

// LoadTraderConfig reads the configuration from a YAML file
func LoadTraderConfig(configFile string) error {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal the config into the traderConfig variable
	if err := viper.Unmarshal(&traderConfig); err != nil {
		return fmt.Errorf("error unmarshalling config: %v", err)
	}
	return nil
}

func GetAppConfig() AppConfig {
	return appConfig
}

func GetTraderConfig() TraderConfig {
	return traderConfig
}
