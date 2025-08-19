package models

import "time"

type Token struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	TokenType             string    `json:"token_type"`
	AccessTokenExpiresAt  time.Time `json:"access_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_expires_at"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token         *Token `json:"token"`
	RequiresOTP   bool   `json:"requires_otp"`
	TransactionID string `json:"transaction_id"`
}

type OTPVerifyInput struct {
	TransactionID string `json:"transaction_id"`
	OTP           string `json:"otp"`
}

type OTPEnrollRequest struct {
	Secret     string `json:"secret"`
	OTPAuthURI string `json:"otp_auth_uri"`
}
