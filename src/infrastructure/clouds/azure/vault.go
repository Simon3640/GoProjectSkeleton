package azure

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"goprojectskeleton/src/infrastructure/config"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

type VaultConfigLoader struct{}

var _ config.Loader = (*VaultConfigLoader)(nil)

func (v *VaultConfigLoader) Load() (*config.Config, error) {

	cfg := &config.Config{}
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
		if err := v.setFieldValue(fieldValue, envValue, defaultValue); err != nil {
			return nil, err
		}
	}

	log.Println("Configuration loaded successfully")
	return cfg, nil
}

func (v *VaultConfigLoader) setFieldValue(field reflect.Value, value string, defaultValue string) error {
	// En Azure Functions, las referencias de Key Vault son resueltas automáticamente
	// por el runtime. Si el valor aún tiene el prefijo, significa que Azure Functions
	// no pudo resolverlo.
	// NUNCA intentamos resolver manualmente las referencias de Key Vault porque:
	// 1. En Azure, el runtime debería resolverlas automáticamente
	// 2. Intentar resolverlas manualmente causa problemas de autenticación
	// 3. Es más seguro usar el valor por defecto si Azure no pudo resolverlo
	if strings.HasPrefix(value, "@Microsoft.KeyVault(SecretUri=") {
		log.Printf("Key Vault reference not resolved: resolving manually")

		secretUri := strings.TrimPrefix(value, "@Microsoft.KeyVault(SecretUri=")
		secretUri = strings.TrimSuffix(secretUri, ")")
		secret, err := v.GetSecret(secretUri)
		if err != nil {
			return fmt.Errorf("failed to get secret '%s' from Key Vault: %w", value, err)
		}
		field.SetString(secret)
		return nil
	}
	field.SetString(value)
	return nil
}

func (v *VaultConfigLoader) GetSecret(secretUri string) (string, error) {
	ctx := context.Background()
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return "", fmt.Errorf("failed to get Azure credentials: %w", err)
	}

	log.Printf("Getting secret '%s' from Key Vault", secretUri)

	vaultUrl, secretName := extractVaultInfo(secretUri)

	client, err := azsecrets.NewClient(vaultUrl, cred, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create Key Vault client: %w", err)
	}

	result, err := client.GetSecret(ctx, secretName, "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to get secret '%s' from Key Vault: %w", secretName, err)
	}

	return *result.Value, nil
}

func NewVaultConfigLoader() *VaultConfigLoader {
	return &VaultConfigLoader{}
}

func extractVaultInfo(uri string) (vaultUrl string, secretName string) {
	parts := strings.Split(uri, "/")
	vaultUrl = strings.Join(parts[:3], "/")
	secretName = parts[4]
	return vaultUrl, secretName
}
