package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Configuration object
type Configuration struct {
	App   AppConfiguration
	Redis RedisConfiguration
}

// LoadConfig returns a configuration object
func LoadConfig() (*Configuration, error) {
	// initialize config & set default
	viper.SetEnvPrefix("US")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("app.host", "localhost:3000")
	viper.SetDefault("app.port", 3000)
	viper.SetDefault("redis.host", "localhost:6379")

	config := new(Configuration)
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
