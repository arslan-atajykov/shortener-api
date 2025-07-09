package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	RedisAddr string
	RedisPass string
	RedisDB   int
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("error loading .env file:", err)
	}

	cfg := &Config{
		Port:      viper.GetString("PORT"),
		DBHost:    viper.GetString("DB_HOST"),
		DBPort:    viper.GetString("DB_PORT"),
		DBUser:    viper.GetString("DB_USER"),
		DBPass:    viper.GetString("DB_PASSWORD"),
		DBName:    viper.GetString("DB_NAME"),
		RedisAddr: viper.GetString("REDIS_ADDR"),
		RedisPass: viper.GetString("REDIS_PASSWORD"),
		RedisDB:   viper.GetInt("REDIS_DB"),
	}
	return cfg
}
