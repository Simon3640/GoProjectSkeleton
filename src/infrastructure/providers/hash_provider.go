package providers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	contracts "gormgoskeleton/src/application/contracts"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashProvider struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

var _ contracts.IHashProvider = (*HashProvider)(nil)

func (hp *HashProvider) HashPassword(password string) (string, error) {
	// creating a salt of variable length
	salt := make([]byte, hp.saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, hp.time, hp.memory, hp.threads, hp.keyLen)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	final := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		hp.memory, hp.time, hp.threads, encodedSalt, encodedHash)

	return final, nil
}

func (hp *HashProvider) VerifyPassword(hashedPassword, password string) (bool, error) {
	// Dividimos el formato
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false, errors.New("hash inválido")
	}

	// Extraer parámetros
	var mem uint32
	var t uint32
	var p uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &mem, &t, &p)
	if err != nil {
		return false, err
	}

	// Extraer salt y hash originales
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Recalcular hash con la contraseña ingresada
	newHash := argon2.IDKey([]byte(password), salt, t, mem, p, uint32(len(hash)))

	return subtle.ConstantTimeCompare(hash, newHash) == 1, nil
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
