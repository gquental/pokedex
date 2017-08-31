package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configuration struct {
	APIEndpoint string `mapstructure:"api_endpoint"`

	Port      string
	DBAddress string `mapstructure:"db_address"`
	Database  string
}

var Config Configuration

func init() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()

	if err != nil {
		panic(fmt.Sprintf("Error loading configuration file: %v", err))
	}

	Config = Configuration{}
	err = v.Unmarshal(&Config)
	if err != nil {
		Config = Configuration{}
	}
}
