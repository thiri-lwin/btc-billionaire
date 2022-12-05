package utilities

import (
	"fmt"

	"github.com/ory/viper"
)

type Config struct {
	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUserName string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")

	viper.AddConfigPath(configPath)

	viper.AutomaticEnv()

	viper.SetConfigType("env")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file, ", err)
		return nil, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to decode into struct, ", err)
		return nil, err
	}

	return &config, nil
}
