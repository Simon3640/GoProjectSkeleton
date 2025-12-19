package settings

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
)

// AppSettings represents the application settings and is used to store the application settings
type AppSettings struct {
	// Application
	AppName         string
	AppEnv          string
	AppPort         string
	AppVersion      string
	AppDescription  string
	AppSupportEmail string
	EnableLog       bool
	DebugLog        bool
	TemplatesPath   string
	AllowOrigins    []string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSL      bool

	// Redis
	RedisHost     string
	RedisPassword string
	RedisDB       int
	RedisTTL      int // in seconds

	// Security
	JWTSecretKey               string
	JWTIssuer                  string
	JWTAudience                string
	JWTAccessTTL               int64 // in seconds
	JWTRefreshTTL              int64 // in seconds
	JWTClockSkew               int64 // in seconds
	OneTimeTokenPasswordTTL    int64 // in minutes
	OneTimeTokenEmailVerifyTTL int64 // in minutes
	FrontendResetPasswordURL   string
	FrontendActivateAccountURL string
	OneTimePasswordTTL         int64 // in minutes
	OneTimePasswordLength      int   // length of the generated one-time password
	LoginMaxAttempts           int   // maximum number of failed login attempts
	LoginAttemptsWindowMinutes int64 // time window in minutes for counting failed attempts

	// Mail
	MailHost         string
	MailPort         int
	MailPassword     string
	MailFrom         string
	MailAuthRequired bool

	// Background Workers
	BackgroundWorkers   int
	BackgroundQueueSize int

	// Observability
	ObservabilityEnabled      bool
	ObservabilityBackend      string
	OTLPEndpoint              string
	ObservabilitySamplingRate string
}

// NewAppSettings creates a new AppSettings instance with default values
func NewAppSettings() *AppSettings {
	return &AppSettings{
		AppName:    "goprojectskeleton",
		AppEnv:     "development",
		AppPort:    "8080",
		AppVersion: "0.0.1",
		EnableLog:  true,
		DebugLog:   true,
	}
}

// Initialize initializes the AppSettings instance with the given values
func (as *AppSettings) Initialize(values map[string]string) *application_errors.ApplicationError {
	asValue := reflect.ValueOf(as).Elem()
	asType := asValue.Type()

	for i := 0; i < asValue.NumField(); i++ {
		field := asType.Field(i)
		fieldValue := asValue.Field(i)
		if !fieldValue.CanSet() {
			continue
		}
		if value, exists := values[field.Name]; exists {
			if err := setFieldValue(fieldValue, value, field.Name); err != nil {
				return application_errors.NewApplicationError(
					status.ApplicationInitializationError,
					messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
					err.Error(),
				)
			}
		}
	}
	return nil
}

func setFieldValue(field reflect.Value, value string, fieldName string) error {
	if value == "" {
		return nil
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		field.SetBool(value == "true" || value == "1" || value == "True" || value == "TRUE")
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s for field %s", value, fieldName)
		}
		field.SetInt(int64(intValue))
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value: %s for field %s", value, fieldName)
		}
		field.SetFloat(floatValue)
	case reflect.Int64:
		int64Value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int64 value: %s for field %s", value, fieldName)
		}
		field.SetInt(int64Value)
	case reflect.Slice:
		values := strings.Split(value, ",")
		slice := reflect.MakeSlice(field.Type(), len(values), len(values))
		for i, v := range values {
			slice.Index(i).SetString(v)
		}
		field.Set(slice)
	default:
		return fmt.Errorf("unsupported field type: %s for field %s", field.Kind(), fieldName)
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
