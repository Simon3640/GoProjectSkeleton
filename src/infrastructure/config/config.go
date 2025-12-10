// Package config provides configuration loading and management functionality.
package config

import (
	"fmt"
	"log"
	"reflect"

	logger "github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// Loader defines the interface for loading configuration from various sources.
type Loader interface {
	Load() (*Config, error)
}

type Config struct {
	// Application
	AppName         string `env:"APP_NAME" envDefault:"0.0.0.1"`
	AppEnv          string `env:"APP_ENV" envDefault:"development"`
	AppPort         string `env:"APP_PORT" envDefault:"8080"`
	AppVersion      string `env:"APP_VERSION" envDefault:"0.0.1"`
	AppDescription  string `env:"APP_DESCRIPTION" envDefault:"Description"`
	AppSupportEmail string `env:"APP_SUPPORT_EMAIL" envDefault:"support@goprojectskeleton.com"`
	EnableLog       string `env:"ENABLE_LOG" envDefault:"true"`
	DebugLog        string `env:"DEBUG_LOG" envDefault:"true"`
	TemplatesPath   string `env:"TEMPLATES_PATH" envDefault:"src/application/shared/templates/"`
	AllowOrigins    string `env:"ALLOW_ORIGINS" envDefault:"*"`

	// Database
	DBHost     string `env:"DB_HOST" envDefault:"goprojectskeleton"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName     string `env:"DB_NAME" envDefault:"goprojectskeleton"`
	DBSSL      string `env:"DB_SSL" envDefault:"false"`

	// Redis
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB       string `env:"REDIS_DB" envDefault:"0"`
	RedisTTL      string `env:"REDIS_TTL" envDefault:"300"`

	// Security
	JWTSecretKey               string `env:"JWT_SECRET_KEY" envDefault:"secret"`
	JWTIssuer                  string `env:"JWT_ISSUER" envDefault:"test-issuer"`
	JWTAudience                string `env:"JWT_AUDIENCE" envDefault:"test-audience"`
	JWTAccessTTL               string `env:"JWT_ACCESS_TTL" envDefault:"3600"`
	JWTRefreshTTL              string `env:"JWT_REFRESH_TTL" envDefault:"86400"`
	JWTClockSkew               string `env:"JWT_CLOCK_SKEW" envDefault:"60"`
	OneTimeTokenPasswordTTL    string `env:"ONE_TIME_TOKEN_TTL" envDefault:"15"`
	OneTimeTokenEmailVerifyTTL string `env:"ONE_TIME_TOKEN_EMAIL_VERIFY_TTL" envDefault:"60"`
	FrontendResetPasswordURL   string `env:"FRONTEND_RESET_PASSWORD_URL" envDefault:"http://localhost:3000/reset-password"`
	FrontendActivateAccountURL string `env:"FRONTEND_ACTIVATE_ACCOUNT_URL" envDefault:"http://localhost:3000/activate-account"`
	OneTimePasswordLength      string `env:"ONE_TIME_PASSWORD_LENGTH" envDefault:"6"`
	OneTimePasswordTTL         string `env:"ONE_TIME_PASSWORD_TTL" envDefault:"10"`
	LoginMaxAttempts           string `env:"LOGIN_MAX_ATTEMPTS" envDefault:"5"`
	LoginAttemptsWindowMinutes string `env:"LOGIN_ATTEMPTS_WINDOW_MINUTES" envDefault:"15"`

	// Mail
	MailHost         string `env:"MAIL_HOST" envDefault:"localhost"`
	MailPort         string `env:"MAIL_PORT" envDefault:"1025"`
	MailPassword     string `env:"MAIL_PASSWORD" envDefault:"password"`
	MailFrom         string `env:"MAIL_FROM" envDefault:"noreply@example.com"`
	MailAuthRequired string `env:"MAIL_AUTH_REQUIRED" envDefault:"true"`
}

func (c *Config) ToMap() map[string]string {
	values := make(map[string]string)
	if c == nil {
		log.Println("Warning: ToMap() called on nil Config")
		return values
	}
	cfgValue := reflect.ValueOf(c).Elem()

	for i := 0; i < cfgValue.NumField(); i++ {
		field := cfgValue.Type().Field(i)
		value := cfgValue.Field(i).String()
		values[field.Name] = value
	}
	return values
}

// NewConfig creates a new configuration instance using the specified loader.
func NewConfig(loader Loader) *Config {
	if loader == nil {
		loader = NewEnvConfigLoader()
	}
	config, err := loader.Load()
	if err != nil {
		fmt.Println("Error loading configuration")
		logger.Logger.Panic("Error loading configuration", err)
		// Si el logger no está habilitado, hacer panic explícitamente
		panic(fmt.Sprintf("Error loading configuration: %v", err))
	}
	return config
}

var ConfigInstance *Config
