package providers

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/settings"
	"gormgoskeleton/src/application/shared/status"

	"golang.org/x/crypto/argon2"
)

type HashProvider struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

var _ contracts_providers.IHashProvider = (*HashProvider)(nil)

func (hp *HashProvider) HashPassword(password string) (string, *application_errors.ApplicationError) {
	// creating a salt of variable settings.AppSettingsInstance.OneTimePasswordLength
	salt := make([]byte, hp.saltLen)
	_, _ = rand.Read(salt)
	// It never return error

	hash := argon2.IDKey([]byte(password), salt, hp.time, hp.memory, hp.threads, hp.keyLen)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	final := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		hp.memory, hp.time, hp.threads, encodedSalt, encodedHash)

	return final, nil
}

func (hp *HashProvider) VerifyPassword(hashedPassword, password string) (bool, *application_errors.ApplicationError) {
	// Dividimos el formato
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false, application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			"invalid hash format")
	}

	// Extraer parámetros
	var mem uint32
	var t uint32
	var p uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &mem, &t, &p)
	if err != nil {
		return false, application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			"failed to parse hash parameters")
	}

	// Extraer salt y hash originales
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			"failed to decode salt")
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			"failed to decode hash")
	}

	// Recalcular hash con la contraseña ingresada
	newHash := argon2.IDKey([]byte(password), salt, t, mem, p, uint32(len(hash)))

	return subtle.ConstantTimeCompare(hash, newHash) == 1, nil
}

func (hp *HashProvider) OneTimeToken() (string, []byte, *application_errors.ApplicationError) {
	// creating a salt of variable settings.AppSettingsInstance.OneTimePasswordLength
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", nil, application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			"failed to generate random bytes")
	}
	token := base64.RawURLEncoding.EncodeToString(salt)
	return token, hp.HashOneTimeToken(token), nil
}

func (hp *HashProvider) HashOneTimeToken(token string) []byte {
	h := sha256.Sum256([]byte(token))
	return h[:]
}

func (hp *HashProvider) ValidateOneTimeToken(hashedToken []byte, token string) bool {
	h := sha256.Sum256([]byte(token))
	return subtle.ConstantTimeCompare(hashedToken, h[:]) == 1
}

func (hp *HashProvider) GenerateOTP() (string, []byte, *application_errors.ApplicationError) {
	digits := "0123456789"
	otp := make([]byte, settings.AppSettingsInstance.OneTimePasswordLength)

	randomBytes := make([]byte, settings.AppSettingsInstance.OneTimePasswordLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", nil, application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			"failed to generate random bytes for OTP",
		)
	}

	for i := 0; i < settings.AppSettingsInstance.OneTimePasswordLength; i++ {
		otp[i] = digits[int(randomBytes[i])%len(digits)]
	}

	return string(otp), hp.HashOneTimeToken(string(otp)), nil
}

// ValidateOTP checks if the provided OTP matches the hashed OTP
func (hp *HashProvider) ValidateOTP(hashedOTP []byte, otp string) bool {
	return hp.ValidateOneTimeToken(hashedOTP, otp)
}

func NewHashProvider() *HashProvider {
	return &HashProvider{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
		saltLen: 16,
	}
}

var HashProviderInstance *HashProvider

func init() {
	HashProviderInstance = NewHashProvider()
}
