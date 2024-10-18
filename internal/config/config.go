package config

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type AppConfig struct {
	IBKR struct {
		Http struct {
			Endpoint string `mapstructure:"Endpoint"`
		} `mapstructure:"Http"`
		Websocket struct {
			Endpoint string `mapstructure:"Endpoint"`
		} `mapstructure:"Websocket"`
		Fields []string `mapstructure:"Fields"`
	} `mapstructure:"IBKR"`
}

type BasicSettings struct {
	Symbol     string `mapstructure:"Symbol"`
	Run        bool   `mapstructure:"Run"`
	RunSeconds int    `mapstructure:"RunSeconds"`
}

type SetUsd struct {
	BidAskTicks    int     `mapstructure:"BidAskTicks"`
	UsdRatioFromOB float64 `mapstructure:"UsdRatioFromOB"`
	MaxUsd         int     `mapstructure:"MaxUsd"`
}

type StrategySettings struct {
	UsdtConvert    bool     `mapstructure:"UsdtConvert"`
	TickerList     []string `mapstructure:"TickerList"`
	ContractIDList []string `mapstructure:"ContractIDList"`
	EntrySide      string   `mapstructure:"EntrySide"`
	TakerFirst     bool     `mapstructure:"TakerFirst"`
	TakerSecond    bool     `mapstructure:"TakerSecond"`
	ProfitRate     struct {
		Long struct {
			ProfitRateLongEntry float64 `mapstructure:"ProfitRateLongEntry"`
			ProfitRateLongExit  float64 `mapstructure:"ProfitRateLongExit"`
		} `mapstructure:"Long"`
		Short struct {
			ProfitRateShortEntry float64 `mapstructure:"ProfitRateShortEntry"`
			ProfitRateShortExit  float64 `mapstructure:"ProfitRateShortExit"`
		} `mapstructure:"Short"`
	} `mapstructure:"ProfitRate"`
	HedgeType struct {
		LongHedgeType  string `mapstructure:"LongHedgeType"`
		ShortHedgeType string `mapstructure:"ShortHedgeType"`
	} `mapstructure:"HedgeType"`
}

type TraderConfig struct {
	Timer struct {
		StartTime string `mapstructure:"StartTime"`
		EndTime   string `mapstructure:"EndTime"`
		StopTime  string `mapstructure:"StopTime"`
	} `mapstructure:"Timer"`
	Inventory struct {
		EntrySide struct {
			MaxNetDelta   float64 `mapstructure:"MaxNetDelta"`
			MaxTotalDelta float64 `mapstructure:"MaxTotalDelta"`
		} `mapstructure:"EntrySide"`
		ProfitRate struct {
			EntryRate        float64 `mapstructure:"entry_rate"`
			ExitRate         float64 `mapstructure:"exit_rate"`
			PunishEntryRate  float64 `mapstructure:"punish_entry_rate"`
			PunishExitRate   float64 `mapstructure:"punish_exit_rate"`
			PunishUnitPerLvg float64 `mapstructure:"punish_unit_per_lvg"`
		} `mapstructure:"ProfitRate"`
	} `mapstructure:"Inventory"`
	HedgeType struct {
		FRHours int `mapstructure:"FR_Hours"`
		FRBar   int `mapstructure:"FR_bar"`
	} `mapstructure:"HedgeType"`
	RiskControl struct {
		RiskCondition struct {
			MaxNetDelta   float64 `mapstructure:"MaxNetDelta"`
			MaxTotalDelta float64 `mapstructure:"MaxTotalDelta"`
		} `mapstructure:"RiskCondition"`
	} `mapstructure:"RiskControl"`
}

type Config struct {
	AppConfig      AppConfig `mapstructure:"AppConfig"`
	StrategyConfig struct {
		BasicSettings    BasicSettings    `mapstructure:"BasicSettings"`
		SetUsd           SetUsd           `mapstructure:"SetUsd"`
		StrategySettings StrategySettings `mapstructure:"StrategySettings"`
	} `mapstructure:"StrategyConfig"`
	TraderConfig TraderConfig `mapstructure:"TraderConfig"`
}

type IBKRClient struct {
	client  *resty.Client
	BaseUrl string
}

var ArbitrageConfig Config

// LoadConfig reads the configuration from a YAML file
func LoadConfig(configFile string) error {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal the config into the ArbitrageConfig variable
	if err := viper.Unmarshal(&ArbitrageConfig); err != nil {
		return fmt.Errorf("error unmarshalling config: %v", err)
	}
	return nil
}

func GetConfig() Config {
	return ArbitrageConfig
}
