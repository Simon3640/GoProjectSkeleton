package dtos

import "time"

type Token struct {
	AccessToken           string    `json:"accessToken"`
	RefreshToken          string    `json:"refreshToken"`
	TokenType             string    `json:"token_type"`
	AccessTokenExpiresAt  time.Time `json:"accessExpiresAt"`
	RefreshTokenExpiresAt time.Time `json:"refresExpiresAt"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token         *Token `json:"token"`
	RequiresOTP   bool   `json:"requiresOtp"`
	TransactionID string `json:"transactionId"`
}

type OTPVerifyInput struct {
	TransactionID string `json:"transactionId"`
	OTP           string `json:"otp"`
}

type OTPEnrollRequest struct {
	Secret     string `json:"secret"`
	OTPAuthURI string `json:"otpAuthUri"`
}
