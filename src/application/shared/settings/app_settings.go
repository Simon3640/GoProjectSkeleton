package settings

import (
	"errors"
	"reflect"
	"strconv"
)

type AppSettings struct {
	// Application
	AppName        string
	AppEnv         string
	AppPort        string
	AppVersion     string
	AppDescription string
	EnableLog      bool
	DebugLog       bool
}

func NewAppSettings() *AppSettings {
	return &AppSettings{
		AppName:    "Gingoskeleton",
		AppEnv:     "development",
		AppPort:    "8080",
		AppVersion: "0.0.1",
		EnableLog:  true,
		DebugLog:   true,
	}
}

func (as *AppSettings) Initialize(values map[string]string) error {
	asValue := reflect.ValueOf(as).Elem()
	asType := asValue.Type()

	for i := 0; i < asValue.NumField(); i++ {
		field := asType.Field(i)
		fieldValue := asValue.Field(i)
		if !fieldValue.CanSet() {
			continue
		}
		if value, exists := values[field.Name]; exists {
			if err := setFieldValue(fieldValue, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	if value == "" {
		return nil
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		field.SetBool(value == "true")
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return errors.New("invalid integer value: " + value)
		}
		field.SetInt(int64(intValue))
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("invalid float value: " + value)
		}
		field.SetFloat(floatValue)
	default:
		return errors.New("unsupported field type")
	}
	return nil
}

func (as *AppSettings) IsDevelopment() bool {
	return as.AppEnv == "development"
}

var AppSettingsInstance *AppSettings

func init() {
	AppSettingsInstance = NewAppSettings()
}
