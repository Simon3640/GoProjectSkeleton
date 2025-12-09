package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	appconfig "github.com/simon3640/goprojectskeleton/src/infrastructure/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManagerConfigLoader struct{}

var _ appconfig.Loader = (*SecretsManagerConfigLoader)(nil)

func (s *SecretsManagerConfigLoader) Load() (*appconfig.Config, error) {
	cfg := &appconfig.Config{}
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

		// Pasar el valor por defecto para usarlo si es necesario
		if err := s.setFieldValue(fieldValue, envValue, defaultValue); err != nil {
			return nil, err
		}
	}

	log.Println("Configuration loaded successfully")
	return cfg, nil
}

func (s *SecretsManagerConfigLoader) setFieldValue(field reflect.Value, value string, defaultValue string) error {
	// En AWS Lambda, si el valor es un ARN de Secrets Manager, lo resolvemos
	// Los ARNs de Secrets Manager tienen el formato: arn:aws:secretsmanager:region:account:secret:name
	if strings.HasPrefix(value, "arn:aws:secretsmanager:") {
		log.Printf("Secrets Manager ARN detected: resolving secret")
		secret, err := s.GetSecret(value)
		if err != nil {
			// Si falla al obtener el secreto, usar el valor por defecto si está disponible
			if defaultValue != "" {
				log.Printf("Failed to get secret '%s', using default value", value)
				field.SetString(defaultValue)
				return nil
			}
			return fmt.Errorf("failed to get secret '%s' from Secrets Manager: %w", value, err)
		}
		field.SetString(secret)
		return nil
	}
	field.SetString(value)
	return nil
}

func (s *SecretsManagerConfigLoader) GetSecret(secretARN string) (string, error) {
	ctx := context.Background()

	// Cargar configuración de AWS usando credenciales por defecto
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	log.Printf("Getting secret from Secrets Manager ARN: %s", secretARN)

	client := secretsmanager.NewFromConfig(awsCfg)

	result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretARN),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get secret from Secrets Manager: %w", err)
	}

	if result.SecretString == nil {
		return "", fmt.Errorf("secret value is nil for ARN: %s", secretARN)
	}

	return *result.SecretString, nil
}

func NewSecretsManagerConfigLoader() *SecretsManagerConfigLoader {
	return &SecretsManagerConfigLoader{}
}
