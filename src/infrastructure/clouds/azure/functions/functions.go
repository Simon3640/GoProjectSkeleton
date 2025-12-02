package main

// Config define la configuración de una función Azure
type Config struct {
	Path           string // e.g., "auth/login"
	HandlerName    string // e.g., "Login"
	Route          string // e.g., "auth/login" (without api/)
	Method         string // "get", "post", "put", "delete"
	AuthLevel      string // "anonymous" or "function"
	NeedsAuth      bool   // needs AuthMiddleware
	NeedsQuery     bool   // needs QueryMiddleware
	HasPathParams  bool   // has path parameters like :id, :identifier, :otp
	PathParamName  string // name of path parameter if any
	PathParamRoute string // route pattern for parsing params
}

// GetAll retorna la lista completa de funciones configuradas
func GetAll() []Config {
	return []Config{
		// Status
		{
			Path:        "status/health_check",
			HandlerName: "GetHealthCheck",
			Route:       "health-check",
			Method:      "get",
			AuthLevel:   "anonymous",
		},
		// Auth routes
		{
			Path:        "auth/login",
			HandlerName: "Login",
			Route:       "auth/login",
			Method:      "post",
			AuthLevel:   "anonymous",
		},
		{
			Path:        "auth/refresh",
			HandlerName: "RefreshAccessToken",
			Route:       "auth/refresh",
			Method:      "post",
			AuthLevel:   "anonymous",
		},
		{
			Path:           "auth/password_reset",
			HandlerName:    "RequestPasswordReset",
			Route:          "auth/password-reset/{identifier}",
			Method:         "get",
			AuthLevel:      "anonymous",
			HasPathParams:  true,
			PathParamName:  "identifier",
			PathParamRoute: "/api/auth/password-reset/:identifier",
		},
		{
			Path:           "auth/login_otp",
			HandlerName:    "LoginOTP",
			Route:          "auth/login-otp/{otp}",
			Method:         "get",
			AuthLevel:      "anonymous",
			HasPathParams:  true,
			PathParamName:  "otp",
			PathParamRoute: "/api/auth/login-otp/:otp",
		},
		// User routes - public
		{
			Path:        "user/create",
			HandlerName: "CreateUser",
			Route:       "user",
			Method:      "post",
			AuthLevel:   "anonymous",
		},
		{
			Path:        "user/create_with_password",
			HandlerName: "CreateUserAndPassword",
			Route:       "user-password",
			Method:      "post",
			AuthLevel:   "anonymous",
		},
		{
			Path:        "user/activate",
			HandlerName: "ActivateUser",
			Route:       "user/activate",
			Method:      "post",
			AuthLevel:   "anonymous",
		},
		// User routes - private
		{
			Path:           "user/get",
			HandlerName:    "GetUser",
			Route:          "user/{id}",
			Method:         "get",
			AuthLevel:      "function",
			NeedsAuth:      true,
			HasPathParams:  true,
			PathParamName:  "id",
			PathParamRoute: "/api/user/:id",
		},
		{
			Path:        "user/get_all",
			HandlerName: "GetAllUser",
			Route:       "user",
			Method:      "get",
			AuthLevel:   "function",
			NeedsAuth:   true,
			NeedsQuery:  true,
		},
		{
			Path:           "user/update",
			HandlerName:    "UpdateUser",
			Route:          "user/{id}",
			Method:         "put",
			AuthLevel:      "function",
			NeedsAuth:      true,
			HasPathParams:  true,
			PathParamName:  "id",
			PathParamRoute: "/api/user/:id",
		},
		{
			Path:           "user/delete",
			HandlerName:    "DeleteUser",
			Route:          "user/{id}",
			Method:         "delete",
			AuthLevel:      "function",
			NeedsAuth:      true,
			HasPathParams:  true,
			PathParamName:  "id",
			PathParamRoute: "/api/user/:id",
		},
		// Password routes
		{
			Path:        "password/create",
			HandlerName: "CreatePassword",
			Route:       "password",
			Method:      "post",
			AuthLevel:   "function",
			NeedsAuth:   true,
		},
		{
			Path:        "password/reset_token",
			HandlerName: "CreatePasswordToken",
			Route:       "password/reset-token",
			Method:      "post",
			AuthLevel:   "anonymous",
		},
	}
}
