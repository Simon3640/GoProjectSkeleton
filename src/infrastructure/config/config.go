package config

import (
	"fmt"
	"reflect"

	logger "gingoskeleton/src/infrastructure/providers"
)

type Config struct {
	// Application
	AppName        string `env:"APP_NAME" envDefault:"Gingoskeleton"`
	AppEnv         string `env:"APP_ENV" envDefault:"development"`
	AppPort        string `env:"APP_PORT" envDefault:"8080"`
	AppVersion     string `env:"APP_VERSION" envDefault:"0.0.1"`
	AppDescription string `env:"APP_DESCRIPTION" envDefault:"Description"`
	EnableLog      string `env:"ENABLE_LOG" envDefault:"true"`
	DebugLog       string `env:"DEBUG_LOG" envDefault:"true"`
}

func (c *Config) ToMap() map[string]string {
	values := make(map[string]string)
	cfgValue := reflect.ValueOf(c).Elem()

	for i := 0; i < cfgValue.NumField(); i++ {
		field := cfgValue.Type().Field(i)
		value := cfgValue.Field(i).String()
		values[field.Name] = value
	}
	return values
}

func NewConfig() *Config {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration")
		logger.Logger.Panic("Error loading configuration", err)
	}
	return config
}

var ConfigInstance *Config

func init() {
	ConfigInstance = NewConfig()
}
