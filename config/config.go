package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`
	Port             string `mapstructure:"PORT"`
	JWTSecret             string `mapstructure:"JWT_SECRET"`
	JWTRefreshSecret      string `mapstructure:"JWT_REFRESH_SECRET"`
	FrontendURL           string `mapstructure:"FRONTEND_URL"`
	AppEnv                string `mapstructure:"APP_ENV"`
	AccessTokenDuration   string `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration  string `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Explicitly bind env vars so Unmarshal works without an .env file
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_SSL_MODE")
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("PORT")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("JWT_REFRESH_SECRET")
	viper.BindEnv("FRONTEND_URL")
	viper.BindEnv("APP_ENV")
	viper.BindEnv("ACCESS_TOKEN_DURATION")
	viper.BindEnv("REFRESH_TOKEN_DURATION")

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error reading config file, using environment variables: %s", err)
		}
	}

	err = viper.Unmarshal(&config)
	return
}
