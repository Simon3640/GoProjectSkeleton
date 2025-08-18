package providers

import (
	"context"
	"time"

	"gormgoskeleton/src/application/contracts"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Secret     []byte
	Issuer     string
	Audience   string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	ClockSkew  time.Duration
}

type JWTProvider struct {
	config Config
}

var _ contracts.IJWTProvider = (*JWTProvider)(nil)

func (jp *JWTProvider) Setup(Secret string, Issuer string, Audience string,
	AccessTTL int64, RefreshTTL int64, ClockSkew int64) {
	jp.config = Config{
		Secret:     []byte(Secret),
		Issuer:     Issuer,
		Audience:   Audience,
		AccessTTL:  time.Duration(AccessTTL) * time.Second,
		RefreshTTL: time.Duration(RefreshTTL) * time.Second,
		ClockSkew:  time.Duration(ClockSkew) * time.Second,
	}
}

func NewJWTProvider() *JWTProvider {
	return &JWTProvider{}
}

func (jp *JWTProvider) GenerateAccessToken(ctx context.Context,
	subject string,
	claimsMap contracts.JWTCLaims) (string, time.Time, error) {
	now := time.Now().Add(jp.config.ClockSkew)
	exp := now.Add(jp.config.AccessTTL)
	claims := jwt.MapClaims{
		"iss": jp.config.Issuer,
		"aud": jp.config.Audience,
		"sub": subject,
		"iat": now.Unix(),
		"nbf": now.Add(-jp.config.ClockSkew).Unix(),
		"exp": exp.Unix(),
	}

	for k, v := range claimsMap {
		claims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jp.config.Secret)
	return signedToken, exp, err
}

func (jp *JWTProvider) GenerateRefreshToken(ctx context.Context,
	subject string) (string, time.Time, error) {
	now := time.Now().Add(jp.config.ClockSkew)
	exp := now.Add(jp.config.RefreshTTL)
	claims := jwt.MapClaims{
		"iss": jp.config.Issuer,
		"aud": jp.config.Audience,
		"sub": subject,
		"iat": now.Unix(),
		"nbf": now.Add(-jp.config.ClockSkew).Unix(),
		"exp": exp.Unix(),
		"typ": "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jp.config.Secret)
	return signedToken, exp, err
}

func (jp *JWTProvider) ParseTokenAndValidate(tokenString string) (contracts.JWTCLaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(jp.config.Secret), nil
	}, jwt.WithAudience(jp.config.Audience), jwt.WithIssuer(jp.config.Issuer), jwt.WithLeeway(jp.config.ClockSkew))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result := make(contracts.JWTCLaims)
		for k, v := range claims {
			result[k] = v
		}
		return result, nil
	} else {
		return nil, jwt.ErrTokenInvalidClaims
	}
}

var JWTProviderInstance *JWTProvider

func init() {
	JWTProviderInstance = NewJWTProvider()
}
