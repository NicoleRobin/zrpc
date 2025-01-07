package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	defaultAddr = ":2345"
)

func GetAddr() string {
	addr := viper.GetString("addr")
	if addr == "" {
		addr = defaultAddr
	}
	return addr
}

func init() {
	viper.SetConfigName("app")  // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
