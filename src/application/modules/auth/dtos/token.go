// Package authdtos contains the DTOs for the auth module.
package authdtos

import "time"

// Token is the DTO for the token
type Token struct {
	AccessToken           string    `json:"accessToken"`
	RefreshToken          string    `json:"refreshToken"`
	TokenType             string    `json:"token_type"`
	AccessTokenExpiresAt  time.Time `json:"accessExpiresAt"`
	RefreshTokenExpiresAt time.Time `json:"refresExpiresAt"`
}

// UserCredentials is the DTO for the user credentials
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is the DTO for the auth response
type AuthResponse struct {
	Token         *Token `json:"token"`
	RequiresOTP   bool   `json:"requiresOtp"`
	TransactionID string `json:"transactionId"`
}

// OTPVerifyInput is the DTO for the OTP verify input
type OTPVerifyInput struct {
	TransactionID string `json:"transactionId"`
	OTP           string `json:"otp"`
}

// OTPEnrollRequest is the DTO for the OTP enroll request
type OTPEnrollRequest struct {
	Secret     string `json:"secret"`
	OTPAuthURI string `json:"otpAuthUri"`
}
