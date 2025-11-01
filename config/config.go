package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort string `mapstructure:"app_port"`

	DBHost     string `mapstructure:"db_host"`
	DBPort     string `mapstructure:"db_port"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBName     string `mapstructure:"db_name"`
	DefaultPage  int `mapstructure:"default_page"`
	DefaultLimit int `mapstructure:"default_limit"`
}

func Load() (config Config, err error) {
	v := viper.New()

	v.SetDefault("app_port", "8080")

	v.BindEnv("app_port", "APP_PORT")
	v.BindEnv("db_host", "DB_HOST")
	v.BindEnv("db_port", "DB_PORT")
	v.BindEnv("db_user", "DB_USER")
	v.BindEnv("db_password", "DB_PASSWORD")
	v.BindEnv("db_name", "DB_NAME")
	v.BindEnv("default_page", "DEFAULT_PAGE")
	v.BindEnv("default_limit", "DEFAULT_LIMIT")

	if err = v.Unmarshal(&config); err != nil {
		return
	}
	log.Printf("AppPort: %s", config.AppPort)
	log.Printf("DBHost: %s", config.DBHost)
	log.Printf("DBUser: %s", config.DBUser)
	log.Printf("DBPort: %s", config.DBPort)
	log.Printf("DBName: %s", config.DBName)
	return
}