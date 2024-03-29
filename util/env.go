package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	PostgresUser         string        `mapstructure:"POSTGRES_USER"`
	PostgresPassword     string        `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDatabase     string        `mapstructure:"POSTGRES_DATABASE"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

// use viper package to read .env file
// return the value of the key
func LoadConfig(path string) (envConfig Env, err error) {
	// viper.AddConfigPath(path)
	// viper.SetConfigName("app")
	// viper.SetConfigType("env")
	viper.SetConfigFile(path)

	viper.AutomaticEnv()

	// Read the config file
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	err = viper.Unmarshal(&envConfig)
	return
}
