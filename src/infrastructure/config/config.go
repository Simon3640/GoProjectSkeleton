package config

import (
	"fmt"
	"reflect"

	logger "gormgoskeleton/src/infrastructure/providers"
)

type Config struct {
	// Application
	AppName        string `env:"APP_NAME" envDefault:"gormgoskeleton"`
	AppEnv         string `env:"APP_ENV" envDefault:"development"`
	AppPort        string `env:"APP_PORT" envDefault:"8080"`
	AppVersion     string `env:"APP_VERSION" envDefault:"0.0.1"`
	AppDescription string `env:"APP_DESCRIPTION" envDefault:"Description"`
	EnableLog      string `env:"ENABLE_LOG" envDefault:"true"`
	DebugLog       string `env:"DEBUG_LOG" envDefault:"true"`

	// Database
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName     string `env:"DB_NAME" envDefault:"gormgoskeleton"`
	DBSSL      string `env:"DB_SSL" envDefault:"false"`

	// Security
	JWTSecretKey  string `env:"JWT_SECRET_KEY" envDefault:"secret"`
	JWTIssuer     string `env:"JWT_ISSUER" envDefault:"test-issuer"`
	JWTAudience   string `env:"JWT_AUDIENCE" envDefault:"test-audience"`
	JWTAccessTTL  string `env:"JWT_ACCESS_TTL" envDefault:"3600"`
	JWTRefreshTTL string `env:"JWT_REFRESH_TTL" envDefault:"86400"`
	JWTClockSkew  string `env:"JWT_CLOCK_SKEW" envDefault:"60"`

	// Mail
	MailHost     string `env:"MAIL_HOST" envDefault:"localhost"`
	MailPort     string `env:"MAIL_PORT" envDefault:"1025"`
	MailPassword string `env:"MAIL_PASSWORD" envDefault:"password"`
	MailFrom     string `env:"MAIL_FROM" envDefault:"noreply@example.com"`
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
