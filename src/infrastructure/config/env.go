package config

import (
	"errors"
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

func LoadConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{}
	cfgValue := reflect.ValueOf(cfg).Elem()

	for i := 0; i < cfgValue.NumField(); i++ {
		field := cfgValue.Type().Field(i)
		envKey := field.Tag.Get("env")
		defaultValue := field.Tag.Get("envDefault")

		if envKey == "" {
			continue
		}

		envValue, exists := os.LookupEnv(envKey)
		if !exists {
			envValue = defaultValue
		}

		fieldValue := cfgValue.Field(i)
		if !fieldValue.CanSet() {
			return nil, errors.New("cannot set value for field: " + field.Name)
		}

		if err := setFieldValue(fieldValue, envValue); err != nil {
			return nil, err
		}
	}

	log.Println("Configuration loaded successfully")
	return cfg, nil
}

func setFieldValue(field reflect.Value, value string) error {
	if value == "" {
		return errors.New("empty value")
	}
	field.SetString(value)
	return nil
}
