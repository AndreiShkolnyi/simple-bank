package util

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")


	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	
	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️ app.env не найден, продолжаем с переменными окружения")
	}

	config = Config{
		DBDriver:            viper.GetString("DB_DRIVER"),
		DBSource:            viper.GetString("DB_SOURCE"),
		ServerAddress:       viper.GetString("SERVER_ADDRESS"),
		TokenSymmetricKey:   viper.GetString("TOKEN_SYMMETRIC_KEY"),
		AccessTokenDuration: viper.GetDuration("ACCESS_TOKEN_DURATION"),
		RefreshTokenDuration: viper.GetDuration("REFRESH_TOKEN_DURATION"),
	}

	return config, nil
}