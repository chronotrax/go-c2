package config

import (
	"fmt"

	"github.com/spf13/viper"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Config struct {
	Debug  bool   `mapstructure:"DEBUG"`
	Port   int    `mapstructure:"PORT"`
	DBName string `mapstructure:"DB_NAME"`
}

func (c Config) String() string {
	return fmt.Sprintf("DEBUG=%t PORT=%d DB_NAME=%s", c.Debug, c.Port, c.DBName)
}

var defaultConfig = Config{
	Debug:  false,
	Port:   80,
	DBName: "db.sqlite",
}

func InitConfig() (Config, error) {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.SetDefault("DEBUG", defaultConfig.Debug)
	viper.SetDefault("PORT", defaultConfig.Port)
	viper.SetDefault("DB_NAME", defaultConfig.DBName)

	err := viper.ReadInConfig()
	if err != nil {
		return defaultConfig, fmt.Errorf("failed to read config file: %w", err)
	}

	// Env variables take precedence over .env file
	viper.AutomaticEnv()

	c := Config{}
	err = viper.Unmarshal(&c)
	if err != nil {
		return defaultConfig, fmt.Errorf("failed to marshal config: %w", err)
	}

	return c, nil
}
