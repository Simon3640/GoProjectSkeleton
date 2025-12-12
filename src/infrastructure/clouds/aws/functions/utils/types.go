// Package utils provides utility functions for generating and deploying AWS Lambda functions.
package utils

// FunctionConfig represents a function configuration from functions.json
type FunctionConfig struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	Handler       string `json:"handler"`
	Route         string `json:"route"`
	Method        string `json:"method"`
	AuthLevel     string `json:"authLevel"`
	NeedsAuth     bool   `json:"needsAuth"`
	NeedsQuery    bool   `json:"needsQuery"`
	HasPathParams bool   `json:"hasPathParams"`
	PathParamName string `json:"pathParamName"`
}

const FunctionsJSONPath = "../functions.json"

// GetHandlerPackage returns the package name for a given handler
func GetHandlerPackage(handlerName string) string {
	handlerPackages := map[string]string{
		// Auth handlers
		"Login":                "authhandlers",
		"RefreshAccessToken":   "authhandlers",
		"RequestPasswordReset": "authhandlers",
		"LoginOTP":             "authhandlers",
		// User handlers
		"CreateUser":            "userhandlers",
		"GetUser":               "userhandlers",
		"UpdateUser":            "userhandlers",
		"DeleteUser":            "userhandlers",
		"GetAllUser":            "userhandlers",
		"CreateUserAndPassword": "userhandlers",
		"ActivateUser":          "userhandlers",
		"ResendWelcomeEmail":    "userhandlers",
		// Password handlers
		"CreatePassword":      "passwordhandlers",
		"CreatePasswordToken": "passwordhandlers",
		// Status handlers
		"GetHealthCheck": "statushandlers",
	}

	if pkg, ok := handlerPackages[handlerName]; ok {
		return pkg
	}
	// Default fallback (should not happen)
	return "handlers"
}

// GetHandlerPackagePath returns the import path for a given handler package
func GetHandlerPackagePath(packageName string) string {
	packagePaths := map[string]string{
		"authhandlers":     "auth",
		"userhandlers":     "user",
		"passwordhandlers": "password",
		"statushandlers":   "status",
	}

	if path, ok := packagePaths[packageName]; ok {
		return path
	}
	return "shared"
}

// GetInitFunction returns the initialization function name for a given handler
func GetInitFunction(handlerName string) string {
	initFunctions := map[string]string{
		// Status handlers
		"GetHealthCheck": "InitializeForStatus",
		// Auth handlers
		"Login":                "InitializeForAuthLogin",
		"RefreshAccessToken":   "InitializeForAuthRefresh",
		"LoginOTP":             "InitializeForAuthLoginOTP",
		"RequestPasswordReset": "InitializeForAuthPasswordReset",
		// User handlers
		"CreateUser":            "InitializeForUser",
		"GetUser":               "InitializeForUser",
		"UpdateUser":            "InitializeForUser",
		"DeleteUser":            "InitializeForUser",
		"ActivateUser":          "InitializeForUser",
		"GetAllUser":            "InitializeForUserWithCache",
		"CreateUserAndPassword": "InitializeForUserWithEmail",
		"ResendWelcomeEmail":    "InitializeForUserWithEmail",
		// Password handlers
		"CreatePassword":      "InitializeForPassword",
		"CreatePasswordToken": "InitializeForPasswordWithEmail",
	}

	if fn, ok := initFunctions[handlerName]; ok {
		return fn
	}
	// Default fallback - initialize everything for safety
	return "InitializeInfrastructure"
}
